package fastrouter

import (
	"strings"
	"sync"

	"github.com/valyala/fasthttp"
)

// paramsKey is the user value key the matched route's params are stored under.
type paramsKey struct{}

// params holds one route's param values. Values point into the request's path buffer instead of copying it,
// so they are valid until the handler returns. A mounted router chains its params to its parent's.
type params struct {
	names  []string
	values []string
	parent *params
}

var paramsPool = sync.Pool{New: func() any { return new(params) }}

// Param returns the value of the named param of the matched route, or "" if there is none.
// The value is valid until the handler returns or changes the request path: use strings.Clone to keep it longer.
func Param(ctx *fasthttp.RequestCtx, name string) string {
	p, _ := ctx.UserValue(paramsKey{}).(*params)
	for ; p != nil; p = p.parent {
		for i, n := range p.names {
			if n == name {
				return p.values[i]
			}
		}
	}

	return ""
}

func setParams(ctx *fasthttp.RequestCtx, names, values []string) {
	p := paramsPool.Get().(*params)
	p.names = names
	p.values = append(p.values[:0], values[:len(names)]...)
	p.parent, _ = ctx.UserValue(paramsKey{}).(*params)
	ctx.SetUserValue(paramsKey{}, p)
}

// detachParams copies the param values out of the path buffer before the path is rewritten.
func detachParams(ctx *fasthttp.RequestCtx) {
	p, _ := ctx.UserValue(paramsKey{}).(*params)
	for ; p != nil; p = p.parent {
		for i, v := range p.values {
			p.values[i] = strings.Clone(v)
		}
	}
}

// Close returns p and its parents to the pool. fasthttp calls it once the request is done.
func (p *params) Close() error {
	for p != nil {
		parent := p.parent
		clear(p.values)
		*p = params{values: p.values[:0]}
		paramsPool.Put(p)
		p = parent
	}

	return nil
}
