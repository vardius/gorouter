package mux

// Route is an handler aware route interface
type Route interface {
	Handler() interface{}
}
