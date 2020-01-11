package gorouter

type route struct {
	handler interface{}
}

func newRoute(h interface{}) *route {
	if h == nil {
		panic("Handler can not be nil.")
	}

	return &route{
		handler: h,
	}
}

func (r *route) Handler() interface{} {
	// returns already cached computed handler
	return r.handler
}
