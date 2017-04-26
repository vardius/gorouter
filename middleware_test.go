package goserver

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestDefaultServeMux(t *testing.T) {
	m := newMiddleware()
	if m.handle(nil) != http.DefaultServeMux {
		t.Error("nil is not DefaultServeMux")
	}
}

func TestHandlerFunc(t *testing.T) {
	fn := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(200)
	})

	m := newMiddleware()
	h := m.handleFunc(fn)

	w := httptest.NewRecorder()

	h.ServeHTTP(w, (*http.Request)(nil))

	if reflect.TypeOf(h) != reflect.TypeOf((http.HandlerFunc)(nil)) {
		t.Error("handleFunc does not construct HandlerFunc")
	}
}

func TestOrders(t *testing.T) {
	m1 := mockMiddlewareWithBody("1")
	m2 := mockMiddlewareWithBody("2")
	m3 := mockMiddlewareWithBody("3")
	fn := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("4"))
	})

	m := newMiddleware(m1, m2, m3)
	h := m.handleFunc(fn)

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
	m1 := mockMiddlewareWithBody("1")
	m2 := mockMiddlewareWithBody("2")
	m3 := mockMiddlewareWithBody("3")
	fn := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("4"))
	})

	m := newMiddleware(m1)
	m = m.append(m2, m3)
	h := m.handleFunc(fn)

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
