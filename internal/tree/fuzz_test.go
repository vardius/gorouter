package tree

import (
	"fmt"
	"math"
	"slices"
	"strings"
	"testing"
)

// refRoute is a route as the reference matcher sees it.
type refRoute struct {
	index    int
	segments []segment
	prefix   bool
}

// refLookup is a deliberately naive matcher: it checks every route against the path,
// then keeps the one the priority rules prefer. FuzzLookup holds the tree to it.
func refLookup(routes []refRoute, path string) (int, []string) {
	var parts []string
	if p := trim(path); p != "" {
		parts = strings.Split(p, "/")
	}

	best, bestRank := -1, []int(nil)
	for _, r := range routes {
		if !refMatches(r, parts) {
			continue
		}
		if rank := refRank(routes, r); best < 0 || slices.Compare(rank, bestRank) < 0 {
			best, bestRank = r.index, rank
		}
	}
	if best < 0 {
		return -1, nil
	}

	var values []string
	for i, s := range routes[slices.IndexFunc(routes, func(r refRoute) bool { return r.index == best })].segments {
		if s.kind != static {
			values = append(values, parts[i])
		}
	}

	return best, values
}

func refMatches(r refRoute, parts []string) bool {
	if len(parts) < len(r.segments) || (!r.prefix && len(parts) != len(r.segments)) {
		return false
	}
	for i, s := range r.segments {
		part := parts[i]
		switch {
		case part == "":
			return false
		case s.kind == static && part != s.text:
			return false
		case s.kind == regexpParam && !s.re.MatchString(part):
			return false
		}
	}

	return true
}

// refRank scores each segment: static first, then regexps in the order they first appeared under their parent, then the wildcard.
func refRank(routes []refRoute, r refRoute) []int {
	rank := make([]int, len(r.segments))
	for k, s := range r.segments {
		switch s.kind {
		case static:
			rank[k] = 0
		case wildcard:
			rank[k] = math.MaxInt
		default:
			var seen []string
			for _, o := range routes {
				if len(o.segments) > k && o.segments[k].kind == regexpParam && sameParent(o.segments, r.segments, k) &&
					!slices.Contains(seen, o.segments[k].text) {
					seen = append(seen, o.segments[k].text)
				}
			}
			rank[k] = 1 + slices.Index(seen, s.text)
		}
	}

	return rank
}

func sameParent(a, b []segment, k int) bool {
	for i := range k {
		if !a[i].sameNode(b[i]) {
			return false
		}
	}

	return true
}

// generator turns fuzz input into patterns and paths from a small alphabet, so routes collide often.
type generator struct{ data []byte }

func (g *generator) next(n int) int {
	if len(g.data) == 0 {
		return 0
	}
	b := g.data[0]
	g.data = g.data[1:]

	return int(b) % n
}

func (g *generator) pattern() (string, bool) {
	kinds := []string{"a", "b", "ab", "c", "{%s}", "{%s:[0-9]+}", "{%s:a|b}", "{%s:[a-c]+}"}
	var sb strings.Builder
	for i := range g.next(4) {
		seg := kinds[g.next(len(kinds))]
		if strings.Contains(seg, "%s") {
			seg = fmt.Sprintf(seg, fmt.Sprintf("p%d", i))
		}
		sb.WriteString("/" + seg)
	}
	if sb.Len() == 0 {
		sb.WriteString("/")
	}

	return sb.String(), g.next(6) == 0
}

func (g *generator) path() string {
	parts := []string{"a", "b", "ab", "c", "1", "12", "x", ""}
	var sb strings.Builder
	for range g.next(5) {
		sb.WriteString("/" + parts[g.next(len(parts))])
	}
	if sb.Len() == 0 || g.next(4) == 0 {
		sb.WriteString("/")
	}

	return sb.String()
}

func FuzzLookup(f *testing.F) {
	f.Add([]byte("\x05\x03\x04\x01\x02\x06\x07\x00\x03\x02\x05\x01\x04\x07"))
	f.Add([]byte("\x0b\x02\x04\x05\x02\x00\x01\x03\x06\x02\x07\x05\x01\x00\x04\x03\x02\x06"))
	f.Add([]byte("\x08\x01\x05\x00\x03\x06\x02\x04\x07\x01\x05\x03\x00\x02\x06\x04"))

	f.Fuzz(func(t *testing.T, data []byte) {
		g := &generator{data: data}

		var tr Tree[int]
		var routes []refRoute
		for i := range g.next(12) + 1 {
			pattern, prefix := g.pattern()
			if !registers(&tr, pattern, prefix, i) {
				continue
			}
			routes = append(routes, refRoute{index: i, segments: mustParse(pattern), prefix: prefix})
		}

		values := make([]string, max(tr.MaxParams(), 1))
		for range 8 {
			path := g.path()
			wantIndex, wantValues := refLookup(routes, path)

			gotIndex, gotValues := -1, []string(nil)
			if r := tr.Lookup("GET", path, values); r != nil {
				gotIndex, gotValues = r.Handler, slices.Clone(values[:len(r.Params)])
			}

			if gotIndex != wantIndex || !slices.Equal(gotValues, wantValues) {
				t.Fatalf("Lookup(%q) = route %d %q, reference says route %d %q\nroutes: %v",
					path, gotIndex, gotValues, wantIndex, wantValues, describe(tr.Routes()))
			}
		}
	})
}

// registers adds the route and reports whether the tree accepted it.
func registers(tr *Tree[int], pattern string, prefix bool, h int) (ok bool) {
	defer func() {
		if recover() != nil {
			ok = false
		}
	}()

	if prefix {
		tr.Mount(pattern, h)
	} else {
		tr.Handle("GET", pattern, h)
	}

	return true
}

func describe(routes []*Route[int]) []string {
	var out []string
	for _, r := range routes {
		if r.Method == "GET" {
			out = append(out, fmt.Sprintf("%d:%s", r.Handler, r.Pattern))
		}
	}

	return out
}
