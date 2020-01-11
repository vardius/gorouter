package mux

// Route is an middleware aware route interface
type Route interface {
	Handler() interface{}
}
