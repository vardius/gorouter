package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func mockMiddleware(body string) MiddlewareFunc {
	fn := func(h interface{}) interface{} {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(body))
			h.(http.Handler).ServeHTTP(w, r)
		})
	}

	return fn
}

func TestOrders(t *testing.T) {
	m1 := mockMiddleware("1")
	m2 := mockMiddleware("2")
	m3 := mockMiddleware("3")
	fn := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("4"))
	})

	m := New(m1, m2, m3)
	h := m.Compose(fn).(http.Handler)

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
	h := m.Compose(fn).(http.Handler)

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
