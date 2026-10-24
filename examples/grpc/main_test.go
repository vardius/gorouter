package main

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io"
	"math/big"
	"net"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"connectrpc.com/grpchealth"
	"github.com/vardius/gorouter/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/encoding/gzip"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

// serve runs h with the example's server settings on a local port, over TLS when cert is set, and returns its address.
// httptest.Server isn't used because it picks the TLS protocols itself, and the server's own choice is under test.
func serve(t *testing.T, h http.Handler, cert *tls.Certificate) string {
	t.Helper()

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	srv := newServer(ln.Addr().String(), h)
	errc := make(chan error, 1)
	if cert != nil {
		srv.TLSConfig = &tls.Config{Certificates: []tls.Certificate{*cert}}
		go func() { errc <- srv.ServeTLS(ln, "", "") }()
	} else {
		go func() { errc <- srv.Serve(ln) }()
	}
	t.Cleanup(func() {
		if err := srv.Close(); err != nil {
			t.Errorf("Close: %v", err)
		}
		if err := <-errc; !errors.Is(err, http.ErrServerClosed) {
			t.Errorf("Serve: %v", err)
		}
	})

	return ln.Addr().String()
}

// healthClient connects grpc-go's client to the health service at addr.
func healthClient(t *testing.T, addr string, creds credentials.TransportCredentials) healthpb.HealthClient {
	t.Helper()

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { conn.Close() })

	return healthpb.NewHealthClient(conn)
}

func TestGRPCAndHTTPOnOnePort(t *testing.T) {
	addr := serve(t, newRouter(time.Minute), nil)
	client := healthClient(t, addr, insecure.NewCredentials())

	resp, err := client.Check(t.Context(), &healthpb.HealthCheckRequest{})
	if err != nil {
		t.Fatal(err)
	}
	if resp.GetStatus() != healthpb.HealthCheckResponse_SERVING {
		t.Errorf("Check = %v, want SERVING", resp.GetStatus())
	}

	// Watch is a server stream; its first message is the current status.
	stream, err := client.Watch(t.Context(), &healthpb.HealthCheckRequest{})
	if err != nil {
		t.Fatal(err)
	}
	if msg, err := stream.Recv(); err != nil || msg.GetStatus() != healthpb.HealthCheckResponse_SERVING {
		t.Errorf("Watch = %v, %v, want SERVING", msg.GetStatus(), err)
	}

	tests := []struct {
		method, target, contentType, body string
		want                              string
	}{
		// The Connect protocol carries the same call as JSON over HTTP/1.1.
		{http.MethodPost, "/grpc.health.v1.Health/Check", "application/json", "{}", `{"status":"SERVING_STATUS_SERVING"}`},
		{http.MethodGet, "/status", "", "", "ok\n"},
	}
	for _, tt := range tests {
		req, err := http.NewRequestWithContext(t.Context(), tt.method, "http://"+addr+tt.target, strings.NewReader(tt.body))
		if err != nil {
			t.Fatal(err)
		}
		if tt.contentType != "" {
			req.Header.Set("Content-Type", tt.contentType)
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			t.Fatal(err)
		}
		if resp.Proto != "HTTP/1.1" || resp.StatusCode != http.StatusOK || string(body) != tt.want {
			t.Errorf("%s %s = %s %d %q, want HTTP/1.1 200 %q", tt.method, tt.target, resp.Proto, resp.StatusCode, body, tt.want)
		}
	}
}

func TestGRPCOverTLS(t *testing.T) {
	cert, roots := selfSigned(t)
	client := healthClient(t, serve(t, newRouter(time.Minute), &cert), credentials.NewClientTLSFromCert(roots, ""))

	resp, err := client.Check(t.Context(), &healthpb.HealthCheckRequest{})
	if err != nil {
		t.Fatal(err)
	}
	if resp.GetStatus() != healthpb.HealthCheckResponse_SERVING {
		t.Errorf("Check = %v, want SERVING", resp.GetStatus())
	}
}

// A gzipped message compresses to a few kilobytes, so the limit has to hold after decompression.
func TestMessageSizeLimit(t *testing.T) {
	client := healthClient(t, serve(t, newRouter(time.Minute), nil), insecure.NewCredentials())

	req := &healthpb.HealthCheckRequest{Service: strings.Repeat("x", maxMessageBytes)}
	_, err := client.Check(t.Context(), req, grpc.UseCompressor(gzip.Name))
	if status.Code(err) != codes.ResourceExhausted {
		t.Errorf("Check with a %d-byte message = %v, want ResourceExhausted", maxMessageBytes, status.Code(err))
	}
}

// Watch sends the current status, then waits for a change that never comes, so only the call's deadline ends it.
// Each client gives up after 5s: a stream still open then means the server didn't end it.
// Over HTTP/1.1, which the Connect protocol streams over too, the server's own timeouts never would.
func TestCallsEndAtTheDeadline(t *testing.T) {
	addr := serve(t, newRouter(500*time.Millisecond), nil)

	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()
	stream, err := healthClient(t, addr, insecure.NewCredentials()).Watch(ctx, &healthpb.HealthCheckRequest{})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := stream.Recv(); err != nil {
		t.Fatal(err)
	}
	if _, err := stream.Recv(); err == nil || ctx.Err() != nil {
		t.Errorf("Watch over HTTP/2: second Recv = %v, want the server to end the stream", err)
	}

	// A Connect stream message is a flags byte and a 4-byte length, then the JSON.
	msg := append([]byte{0, 0, 0, 0, 2}, "{}"...)
	req, err := http.NewRequestWithContext(t.Context(), http.MethodPost, "http://"+addr+"/grpc.health.v1.Health/Watch", bytes.NewReader(msg))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/connect+json")
	resp, err := (&http.Client{Timeout: 5 * time.Second}).Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if _, err := io.ReadAll(resp.Body); err != nil {
		t.Errorf("Watch over %s: %v, want the server to end the stream", resp.Proto, err)
	}
}

func TestMountHidesTheMethod(t *testing.T) {
	path, health := grpchealth.NewHandler(grpchealth.NewStaticChecker())
	var seen atomic.Value
	r := gorouter.New()
	r.Mount(path, http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		seen.Store(req.URL.Path)
		health.ServeHTTP(w, req)
	}))

	_, err := healthClient(t, serve(t, r, nil), insecure.NewCredentials()).Check(t.Context(), &healthpb.HealthCheckRequest{})
	if got, _ := seen.Load().(string); got != "/Check" || status.Code(err) != codes.Unimplemented {
		t.Errorf("Check through Mount: handler saw %q and the call returned %v, want /Check and Unimplemented", got, err)
	}
}

// selfSigned returns a certificate for 127.0.0.1 and a pool that trusts it.
func selfSigned(t *testing.T) (tls.Certificate, *x509.CertPool) {
	t.Helper()

	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	tmpl := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		IPAddresses:  []net.IP{net.IPv4(127, 0, 0, 1)},
		NotAfter:     time.Now().Add(time.Hour),
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
	if err != nil {
		t.Fatal(err)
	}
	leaf, err := x509.ParseCertificate(der)
	if err != nil {
		t.Fatal(err)
	}
	roots := x509.NewCertPool()
	roots.AddCert(leaf)

	return tls.Certificate{Certificate: [][]byte{der}, PrivateKey: key, Leaf: leaf}, roots
}
