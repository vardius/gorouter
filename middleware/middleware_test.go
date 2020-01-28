package middleware

import (
	"testing"
)

func TestMiddleware_WithPriority(t *testing.T) {
	type test struct {
		name       string
		middleware Middleware
		priority       uint
	}
	tests := []test{
		{"Zero", mockMiddleware("Zero"), 0},
		{"Positive", mockMiddleware("Positive"), 1},
		{"Positive Large", mockMiddleware("Positive Large"), 999},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := WithPriority(tt.middleware, tt.priority)
			if got := m.Priority(); got != tt.priority {
				t.Errorf("Priority() = %v, want %v", got, tt.priority)
			}
		})
	}
}
