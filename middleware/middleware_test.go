package middleware

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func mockMiddleware(body string) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(body))
			h.ServeHTTP(w, r)
		})
	}
}

func TestDefaultServeMux(t *testing.T) {
	m := New()
	if m.Handle(nil) != http.DefaultServeMux {
		t.Error("nil is not DefaultServeMux")
	}
}

func TestHandlerFunc(t *testing.T) {
	fn := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(200)
	})

	m := New()
	h := m.HandleFunc(fn)

	w := httptest.NewRecorder()

	h.ServeHTTP(w, (*http.Request)(nil))

	if reflect.TypeOf(h) != reflect.TypeOf((http.HandlerFunc)(nil)) {
		t.Error("handleFunc does not construct HandlerFunc")
	}
}

func TestOrders(t *testing.T) {
	m1 := mockMiddleware("1")
	m2 := mockMiddleware("2")
	m3 := mockMiddleware("3")
	fn := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("4"))
	})

	m := New(m1, m2, m3)
	h := m.HandleFunc(fn)

	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	h.ServeHTTP(w, r)

	if w.Body.String() != "1234" {
		t.Error("The order is incorrect")
	}
}

func TestAppend(t *testing.T) {
	m1 := mockMiddleware("1")
	m2 := mockMiddleware("2")
	m3 := mockMiddleware("3")
	fn := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("4"))
	})

	m := New(m1)
	m = m.Append(m2, m3)
	h := m.HandleFunc(fn)

	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	h.ServeHTTP(w, r)

	if w.Body.String() != "1234" {
		t.Errorf("The order is incorrect expected: 1234 actual: %s", w.Body.String())
	}
}
