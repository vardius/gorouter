---
id: https
title: HTTPS
sidebar_label: HTTPS
---

## Autocert
<!--DOCUSAURUS_CODE_TABS-->
<!--net/http-->
```go
package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"
	"log"

	"github.com/caarlos0/env"
	"github.com/vardius/gorouter/v4"
	"golang.org/x/crypto/acme/autocert"
)

type config struct {
	Host         string   `env:"HOST"              envDefault:"0.0.0.0"`
	Port         int      `env:"PORT"              envDefault:"3000"`
	CertDirCache string   `env:"CERT_DIR_CACHE"`
}

var localHostAddresses = map[string]bool{
	"0.0.0.0":   true,
	"localhost": true,
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello!\n")
}

func main() {
	cfg := config{}
	env.Parse(&cfg)

	router := gorouter.New()
	router.GET("/", http.HandlerFunc(index))	

	srv := setupServer(&cfg, router)

	log.Fatal(srv.ListenAndServeTLS("", ""))
}

func setupServer(cfg *config, router gorouter.Router) *http.Server {
	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      router,
	}

	// for localhost do not use autocert
	if localHostAddresses[cfg.Host] {
		return srv
	}

	certManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(cfg.Host),
		Cache:      autocert.DirCache(cfg.CertDirCache),
	}

	tlsConfig := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
		GetCertificate: certManager.GetCertificate,
	}

	srv.TLSConfig = tlsConfig
	srv.TLSNextProto = make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0)

	return srv
}
```
<!--valyala/fasthttp-->
```go
Example coming soon...
```
<!--END_DOCUSAURUS_CODE_TABS-->