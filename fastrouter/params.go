package fastrouter

import (
	"sync"
	"unsafe"

	"github.com/valyala/fasthttp"
)

// paramsKey is the user value key the request's params are stored under.
type paramsKey struct{}

// params holds the param values of the routes a request matched: a mounted router appends its route's
// after the mount's. Values point into a copy of the path held in buf, so rewriting the request path
// leaves them intact.
type params struct {
	names  []string
	values []string
	buf    []byte
}

var paramsPool = sync.Pool{New: func() any { return new(params) }}

// Param returns the value of the named param of the matched route, or "" if there is none.
// Inside a mounted router, its route's params shadow the mount's params of the same name.
// The value is valid until the handler returns: use strings.Clone to keep it longer.
func Param(ctx *fasthttp.RequestCtx, name string) string {
	p, _ := ctx.UserValue(paramsKey{}).(*params)
	if p == nil {
		return ""
	}

	for i := len(p.names) - 1; i >= 0; i-- {
		if p.names[i] == name {
			return p.values[i]
		}
	}

	return ""
}

// setParams adds a matched route's param values, which point into path, to the request's params.
func setParams(ctx *fasthttp.RequestCtx, path string, names, values []string) {
	p, _ := ctx.UserValue(paramsKey{}).(*params)
	if p == nil {
		p = paramsPool.Get().(*params)
		ctx.SetUserValue(paramsKey{}, p)
	}

	// One copy of the path for all the values. If the append moves buf, values already taken
	// keep the old array alive and unchanged.
	start := len(p.buf)
	p.buf = append(p.buf, path...)
	own := bytesToString(p.buf[start:])

	p.names = append(p.names, names...)
	for _, v := range values[:len(names)] {
		i := offset(path, v)
		p.values = append(p.values, own[i:i+len(v)])
	}
}

// offset returns where sub starts in s. sub must be a substring of s.
func offset(s, sub string) int {
	return int(uintptr(unsafe.Pointer(unsafe.StringData(sub))) - uintptr(unsafe.Pointer(unsafe.StringData(s))))
}

// paramsMark records how many params a request has, so a mount can drop the ones its handler added.
type paramsMark struct {
	params, bytes int
}

func markParams(ctx *fasthttp.RequestCtx) paramsMark {
	p, _ := ctx.UserValue(paramsKey{}).(*params)
	if p == nil {
		return paramsMark{}
	}

	return paramsMark{len(p.names), len(p.buf)}
}

// resetParams drops the params added since m was taken.
func resetParams(ctx *fasthttp.RequestCtx, m paramsMark) {
	p, _ := ctx.UserValue(paramsKey{}).(*params)
	if p == nil {
		return // the handler reset the user values
	}

	n := min(m.params, len(p.names))
	clear(p.values[n:])
	p.names, p.values = p.names[:n], p.values[:n]
	p.buf = p.buf[:min(m.bytes, len(p.buf))]
}

// Close returns p to the pool. fasthttp calls it once the request is done.
func (p *params) Close() error {
	clear(p.values)
	p.names, p.values, p.buf = p.names[:0], p.values[:0], p.buf[:0]
	paramsPool.Put(p)

	return nil
}
