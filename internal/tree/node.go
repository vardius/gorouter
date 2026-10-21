package tree

import (
	"regexp"
	"strings"
)

// node is one segment position in a method's tree.
// Children are tried in priority order: static, then regexp params in registration order, then the wildcard.
type node[H any] struct {
	kind kind
	text string
	re   *regexp.Regexp

	param  int // where a param node writes its value
	params int // number of params from the root down to this node, itself included

	indices  string // the first byte of each static child, in the order of statics
	statics  []*node[H]
	regexps  []*node[H]
	wildcard *node[H]

	route  *Route[H]
	prefix bool // route is a mount, so it also matches every path below this node
}

// child returns the child for s, adding it if needed.
func (n *node[H]) child(s segment) *node[H] {
	switch s.kind {
	case static:
		for _, c := range n.statics {
			if c.text == s.text {
				return c
			}
		}
		c := n.newChild(s)
		n.indices += s.text[:1]
		n.statics = append(n.statics, c)
		return c
	case regexpParam:
		for _, c := range n.regexps {
			if c.text == s.text {
				return c
			}
		}
		c := n.newChild(s)
		n.regexps = append(n.regexps, c)
		return c
	default:
		if n.wildcard == nil {
			n.wildcard = n.newChild(s)
		}
		return n.wildcard
	}
}

func (n *node[H]) newChild(s segment) *node[H] {
	c := &node[H]{kind: s.kind, text: s.text, re: s.re, params: n.params}
	if s.kind != static {
		c.param = n.params
		c.params++
	}

	return c
}

func (n *node[H]) isLeaf() bool {
	return len(n.statics) == 0 && len(n.regexps) == 0 && n.wildcard == nil
}

// isEmpty reports whether no route lives at or below n.
func (n *node[H]) isEmpty() bool {
	return n.route == nil && n.isLeaf()
}

// descend is called once n has matched a segment: rest is the path after it, more says whether there was one.
func (n *node[H]) descend(rest string, more bool, values []string) *Route[H] {
	if n.prefix || !more {
		return n.route
	}

	return n.match(rest, values)
}

// match finds the route for path, which has no leading slash, among n's children.
// A failed branch falls back to the next candidate, so a static miss deeper down can still match a param.
func (n *node[H]) match(path string, values []string) *Route[H] {
	if path == "" || path[0] == '/' {
		return nil // an empty segment never matches
	}

	// A static child is compared against the start of path, so finding where the segment ends is left to params.
	for i := range len(n.indices) {
		c := n.statics[i]
		if n.indices[i] != path[0] || !strings.HasPrefix(path, c.text) {
			continue
		}
		rest := path[len(c.text):]
		if rest != "" && rest[0] != '/' {
			continue // c.text is only a prefix of this segment
		}
		if r := c.descend(strings.TrimPrefix(rest, "/"), rest != "", values); r != nil {
			return r
		}
		break // no other static child can match the same segment
	}

	if len(n.regexps) == 0 && n.wildcard == nil {
		return nil
	}

	seg, rest, more := strings.Cut(path, "/")

	for _, c := range n.regexps {
		if c.re.MatchString(seg) {
			values[c.param] = seg
			if r := c.descend(rest, more, values); r != nil {
				return r
			}
		}
	}

	if c := n.wildcard; c != nil {
		values[c.param] = seg
		return c.descend(rest, more, values)
	}

	return nil
}
