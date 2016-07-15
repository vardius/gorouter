package goapi

import (
	"errors"
	"net/http"
	"strings"
)

type Context struct {
	Params map[string]string
	Route  Route
}

func fromRequest(r *route, req *http.Request) (*Context, error) {
	paths := strings.Split(strings.Trim(req.URL.Path, "/"), "/")
	route, params := r.getRoute(paths)

	if route != nil {
		return &Context{params, route}, nil
	}

	return nil, errors.New("goapi: error while creating context")
}
