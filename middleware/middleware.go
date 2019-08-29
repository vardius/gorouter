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

// Append appends handlers to middlewares
func (m Middleware) Append(fs ...MiddlewareFunc) Middleware {
	for _, f := range fs {
		m = append(m, f)
	}

	return m
}

// Merge merges another middlewares
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
