package middleware

import (
	"testing"
)

type mockWrapper struct{}

func (*mockWrapper) Wrap(h Handler) Handler {
	return h
}

func TestNew(t *testing.T) {
	type args struct {
		w        Wrapper
		priority uint
	}
	type test struct {
		name string
		args args
	}
	tests := []test{
		{"From Wrapper", args{&mockWrapper{}, 0}},
		{"From WrapperFunc", args{WrapperFunc(func(h Handler) Handler { return func() {} }), 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			panicked := false
			defer func() {
				if rcv := recover(); rcv != nil {
					panicked = true
				}
			}()

			got := New(tt.args.w, tt.args.priority)

			if panicked {
				t.Errorf("Panic: New() = %v", got)
			}
		})
	}
}

func TestMiddleware_Priority(t *testing.T) {
	type test struct {
		name       string
		middleware Middleware
		want       uint
	}
	tests := []test{
		{"Zero", mockMiddleware("TestMiddleware_Priority 1", 0), 0},
		{"Positive", mockMiddleware("TestMiddleware_Priority 1", 1), 1},
		{"Positive Large", mockMiddleware("TestMiddleware_Priority 1", 999), 999},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Middleware{
				wrapper:  tt.middleware.wrapper,
				priority: tt.middleware.priority,
			}
			if got := m.Priority(); got != tt.want {
				t.Errorf("Priority() = %v, want %v", got, tt.want)
			}
		})
	}
}
