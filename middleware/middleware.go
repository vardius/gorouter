package middleware

// MiddlewareFunc is a middleware function type.
// Long story - short: it is a handler wrapper
type MiddlewareFunc func(interface{}) interface{}

// Middleware is a slice of handler functions
type Middleware []MiddlewareFunc

// New provides new middleware
func New(fs ...MiddlewareFunc) Middleware {
	return fs
}

// Append appends handlers to middleware
func (m Middleware) Append(fs ...MiddlewareFunc) Middleware {
	return m.Merge(fs)
}

// Merge merges another middleware
func (m Middleware) Merge(n Middleware) Middleware {
	return append(m, n...)
}

// Compose returns middleware composed to single MiddlewareFunc
func (m Middleware) Compose(h interface{}) interface{} {
	if h == nil {
		return nil
	}

	for i := range m {
		h = m[len(m)-1-i](h)
	}

	return h
}
