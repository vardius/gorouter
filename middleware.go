package goapi

import (
	"net/http"
	"sort"
)

type (
	middlewares []*middleware
	middleware  struct {
		path     string
		priority int
		handler  MiddlewareFunc
	}
	MiddlewareFunc func(*http.Request, *Context) Error
)

func (m middlewares) Len() int           { return len(m) }
func (m middlewares) Swap(i, j int)      { m[i], m[j] = m[j], m[i] }
func (m middlewares) Less(i, j int) bool { return m[i].priority < m[j].priority }

func sortByPriority(m middlewares) {
	sort.Sort(middlewares(m))
}
