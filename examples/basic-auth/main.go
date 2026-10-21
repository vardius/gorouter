// Command basic-auth protects everything below /admin with HTTP basic authentication.
// Set ADMIN_USER and ADMIN_PASSWORD before starting it.
package main

import (
	"crypto/sha256"
	"crypto/subtle"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/vardius/gorouter/v5"
)

func newRouter(user, password string) *gorouter.Router {
	r := gorouter.New()
	r.Use("/admin", basicAuth("admin", user, password))

	r.GET("/", index)
	r.GET("/admin", dashboard)
	r.GET("/admin/users/{name}", adminUser)

	return r
}

func index(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintln(w, "public page")
}

func dashboard(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintln(w, "admin dashboard")
}

func adminUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "managing %s\n", r.PathValue("name"))
}

// basicAuth lets a request through only with the given credentials.
// Comparing fixed-size hashes in constant time keeps the check from leaking their length or content.
func basicAuth(realm, user, password string) gorouter.Middleware {
	wantUser, wantPassword := sha256.Sum256([]byte(user)), sha256.Sum256([]byte(password))

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			u, p, ok := r.BasicAuth()
			gotUser, gotPassword := sha256.Sum256([]byte(u)), sha256.Sum256([]byte(p))
			userOK := subtle.ConstantTimeCompare(gotUser[:], wantUser[:])
			passwordOK := subtle.ConstantTimeCompare(gotPassword[:], wantPassword[:])
			if !ok || userOK&passwordOK != 1 {
				w.Header().Set("WWW-Authenticate", `Basic realm="`+realm+`", charset="UTF-8"`)
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func main() {
	user, password := os.Getenv("ADMIN_USER"), os.Getenv("ADMIN_PASSWORD")
	if user == "" || password == "" {
		log.Fatal("set ADMIN_USER and ADMIN_PASSWORD")
	}

	srv := &http.Server{Addr: ":8080", Handler: newRouter(user, password), ReadHeaderTimeout: 5 * time.Second}
	log.Fatal(srv.ListenAndServe())
}
