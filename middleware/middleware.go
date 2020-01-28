package middleware

// Handler represents wrapped function
type Handler interface{}

// Middleware wraps Handler
type Middleware interface {
	// Wrap Handler with middleware
	Wrap(Handler) Handler
	// Priority provides a value for sorting Collection, lower values come first
	Priority() uint
}

// WrapperFunc is an adapter to allow the use of
// handler wrapper functions as middleware functions.
type WrapperFunc func(Handler) Handler

// Wrap implements Wrapper interface
func (f WrapperFunc) Wrap(h Handler) Handler {
	return f(h)
}

// Priority provides a value for sorting Collection, lower values come first
func (f WrapperFunc) Priority() (priority uint) {
	return
}

// Middleware is a slice of handler wrappers functions
type sortableMiddleware struct {
	Middleware
	priority uint
}

// Priority provides a value for sorting Collection, lower values come first
func (m *sortableMiddleware) Priority() uint {
	return m.priority
}

// WithPriority provides new Middleware with priority
func WithPriority(middleware Middleware, priority uint) Middleware {
	return &sortableMiddleware{
		Middleware: middleware,
		priority: priority,
	}
}
