package tree

import (
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strings"
)

type kind uint8

const (
	static kind = iota
	regexpParam
	wildcard
)

// segment is one slash-separated part of a route pattern.
type segment struct {
	kind kind
	text string // the literal for a static segment, the expression for a regexp param
	name string // the param name; empty for a static segment
	re   *regexp.Regexp
}

// sameNode reports whether s and o live in the same tree node.
// Param names don't count: they belong to the route, not to the node.
func (s segment) sameNode(o segment) bool {
	return s.kind == o.kind && s.text == o.text
}

// parse splits a pattern such as /users/{id:[0-9]+}/posts into its segments.
func parse(pattern string) ([]segment, error) {
	if !strings.HasPrefix(pattern, "/") {
		return nil, errors.New("must start with /")
	}
	if strings.Contains(pattern, "//") {
		return nil, errors.New("has an empty segment")
	}

	path := trim(pattern)
	if path == "" {
		return nil, nil
	}

	var segments []segment
	for part := range strings.SplitSeq(path, "/") {
		s, err := parseSegment(part)
		if err != nil {
			return nil, err
		}
		if s.name != "" && slices.ContainsFunc(segments, func(o segment) bool { return o.name == s.name }) {
			return nil, fmt.Errorf("repeats param {%s}", s.name)
		}
		segments = append(segments, s)
	}

	return segments, nil
}

func parseSegment(part string) (segment, error) {
	if part[0] != '{' {
		if strings.ContainsAny(part, "{}") {
			return segment{}, fmt.Errorf("segment %q mixes text and a param", part)
		}
		return segment{kind: static, text: part}, nil
	}

	if part[len(part)-1] != '}' {
		return segment{}, fmt.Errorf("param %q is not closed (a regexp can't contain /)", part)
	}

	name, expr, hasExpr := strings.Cut(part[1:len(part)-1], ":")
	switch {
	case name == "" || strings.ContainsAny(name, "{}"):
		return segment{}, fmt.Errorf("param %q needs a name", part)
	case !hasExpr:
		return segment{kind: wildcard, name: name}, nil
	case expr == "":
		return segment{}, fmt.Errorf("param %q has an empty regexp", part)
	}

	// Anchored, so the expression has to match the whole segment.
	re, err := regexp.Compile("^(?:" + expr + ")$")
	if err != nil {
		return segment{}, fmt.Errorf("param %q: %w", part, err)
	}

	return segment{kind: regexpParam, text: expr, name: name, re: re}, nil
}

// trim drops one leading and one trailing slash: "/a/b/" becomes "a/b", "/" becomes "".
// "//" becomes "/", an empty segment, so it matches nothing.
func trim(path string) string {
	if path != "" && path[0] == '/' {
		path = path[1:]
	}
	if len(path) > 1 && path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}

	return path
}

// Depth returns the number of segments in pattern. It panics if the pattern is invalid.
func Depth(pattern string) int {
	return len(mustParse(pattern))
}

// ParamNames returns the names of pattern's params, in order. It panics if the pattern is invalid.
func ParamNames(pattern string) []string {
	var names []string
	for _, s := range mustParse(pattern) {
		if s.kind != static {
			names = append(names, s.name)
		}
	}

	return names
}

// StripSegments drops the first n segments of a request path,
// so StripSegments("/v1/users/7", 1) is "/users/7". The result is never empty.
func StripSegments(path string, n int) string {
	rest := path
	for range n {
		rest = strings.TrimPrefix(rest, "/")
		i := strings.IndexByte(rest, '/')
		if i < 0 {
			return "/"
		}
		rest = rest[i:]
	}

	if rest == "" {
		return "/"
	}

	return rest
}

func mustParse(pattern string) []segment {
	segments, err := parse(pattern)
	if err != nil {
		panic(fmt.Sprintf("gorouter: pattern %q %v", pattern, err))
	}

	return segments
}
