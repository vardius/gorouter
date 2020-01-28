package middleware

// Handler represents wrapped function
type Handler interface{}

// Wrapper wraps Handler
type Wrapper interface {
	// Wrap Handler with middleware
	Wrap(Handler) Handler
}

// Sortable allows Collection to be sorted by priority
type Sortable interface {
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

// Middleware is a slice of handler wrappers functions
type Middleware struct {
	wrapper  Wrapper
	priority uint
}

// Wrap Handler with middleware
func (m Middleware) Wrap(h Handler) Handler {
	return m.wrapper.Wrap(h)
}

// Priority provides a value for sorting Collection, lower values come first
func (m Middleware) Priority() uint {
	return m.priority
}

// New provides new Middleware
func New(w Wrapper, priority uint) Middleware {
	return Middleware{
		wrapper:  w,
		priority: priority,
	}
}
