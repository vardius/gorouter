package mux

import (
	"testing"

	"github.com/vardius/gorouter/v4/context"
)

func TestNewNode(t *testing.T) {
	node := NewNode("lang", 0)

	switch node := node.(type) {
	case *regexpNode:
		t.Fatalf("Expecting: *mux.staticNode. Wrong node type: %T\n", node)
	case *wildcardNode:
		t.Fatalf("Expecting: *mux.staticNode. Wrong node type: %T\n", node)
	}

	node = NewNode("{lang:en|pl}", 0)

	switch node := node.(type) {
	case *staticNode:
		t.Fatalf("Expecting: *mux.regexpNode. Wrong node type: %T\n", node)
	case *wildcardNode:
		t.Fatalf("Expecting: *mux.regexpNode. Wrong node type: %T\n", node)
	}

	node = NewNode("{lang}", 0)

	switch node := node.(type) {
	case *staticNode:
		t.Fatalf("Expecting: *mux.wildcardNode. Wrong node type: %T\n", node)
	case *regexpNode:
		t.Fatalf("Expecting: *mux.wildcardNode. Wrong node type: %T\n", node)
	}
}

type mockroute struct {
	handler interface{}
}

func newRoute(h interface{}) *mockroute {
	if h == nil {
		panic("Handler can not be nil.")
	}

	return &mockroute{
		handler: h,
	}
}

func (r *mockroute) Handler() interface{} {
	return r.handler
}

func TestStaticNodeMatchRoute(t *testing.T) {
	homeRoute := newRoute("testhomeroute")
	searchRoute := newRoute("testsearchroute")

	homeSkipSubpath := staticNode{name: "home", route: homeRoute, maxParamsSize: 3}
	homeSkipSubpath.SkipSubPath()

	home := staticNode{name: "home", route: nil, maxParamsSize: 3}
	home.WithRoute(homeRoute)

	search := NewNode("search", home.MaxParamsSize())
	search.WithRoute(searchRoute)

	home.WithChildren(home.Tree().withNode(search).sort())
	home.WithChildren(home.Tree().Compile())

	tests := []struct {
		name           string
		node           staticNode
		path           string
		expectedRoute  Route
		expectedParams context.Params
	}{
		{
			name:           "Exact Match",
			node:           home,
			path:           "home",
			expectedRoute:  homeRoute,
			expectedParams: make(context.Params, 3),
		},
		{
			name:           "Exact Match with Skip SubPath",
			node:           homeSkipSubpath,
			path:           "home/about",
			expectedRoute:  homeRoute,
			expectedParams: make(context.Params, 3),
		},
		{
			name:           "Match with SubPath",
			node:           home,
			path:           "home/search",
			expectedRoute:  searchRoute,
			expectedParams: make(context.Params, 3),
		},
		{
			name:           "No Match",
			node:           home,
			path:           "about",
			expectedRoute:  nil,
			expectedParams: nil,
		},
		{
			name:           "No Match without Slash",
			node:           home,
			path:           "homeabout",
			expectedRoute:  nil,
			expectedParams: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			route, params := tt.node.MatchRoute(tt.path)
			if route != tt.expectedRoute {
				t.Errorf("%s: expected route %v, got %v", tt.name, tt.expectedRoute, route)
			}
			if len(params) != len(tt.expectedParams) {
				t.Errorf("%s: expected params %v, got %v", tt.name, tt.expectedParams, params)
			}
		})
	}
}
