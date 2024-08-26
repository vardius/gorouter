package mux

import (
	"reflect"
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

func newMockRoute(h interface{}) *mockroute {
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
	paramSize := 3
	homeRoute := newMockRoute("testhomeroute")
	searchRoute := newMockRoute("testsearchroute")

	homeSkipSubpath := staticNode{name: "home", route: homeRoute, maxParamsSize: uint8(paramSize)}
	homeSkipSubpath.SkipSubPath()

	home := staticNode{name: "home", route: nil, maxParamsSize: uint8(paramSize)}
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
			expectedParams: make(context.Params, paramSize),
		},
		{
			name:           "Exact Match with Skip SubPath",
			node:           homeSkipSubpath,
			path:           "home/about",
			expectedRoute:  homeRoute,
			expectedParams: make(context.Params, paramSize),
		},
		{
			name:           "Match with SubPath",
			node:           home,
			path:           "home/search",
			expectedRoute:  searchRoute,
			expectedParams: make(context.Params, paramSize),
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
			if !reflect.DeepEqual(params, tt.expectedParams) {
				t.Errorf("%s: expected params %v, got %v", tt.name, tt.expectedParams, params)
			}
		})
	}
}

func TestWildcardNodeMatchRoute(t *testing.T) {
	paramSize := 3
	productRoute := newMockRoute("testproductroute")
	itemRoute := newMockRoute("testitemroute")
	viewRoute := newMockRoute("testviewroute")
	params := make(context.Params, paramSize)

	product := staticNode{name: "product", route: nil, maxParamsSize: uint8(paramSize)}
	product.WithRoute(productRoute)

	// Create a wildcard node for "{item}" under "product"
	item := NewNode("{item}", product.MaxParamsSize())
	item.WithRoute(itemRoute)

	view := NewNode("view", product.MaxParamsSize()+1)
	view.WithRoute(viewRoute)

	// Build the tree structure
	product.WithChildren(product.Tree().withNode(item).sort())
	item.WithChildren(item.Tree().withNode(view).sort())
	product.WithChildren(product.Tree().Compile())

	// Create a static node for "product" with skip subpath functionality
	productSkipSubpath := staticNode{name: "product", route: nil, maxParamsSize: uint8(paramSize)}
	productSkipSubpath.WithRoute(productRoute)

	itemSkipSubpath := NewNode("{item}", productSkipSubpath.MaxParamsSize()) // wildcardNode
	itemSkipSubpath.WithRoute(itemRoute)
	itemSkipSubpath.SkipSubPath()

	// Build the tree structure
	productSkipSubpath.WithChildren(productSkipSubpath.Tree().withNode(itemSkipSubpath).sort())
	productSkipSubpath.WithChildren(productSkipSubpath.Tree().Compile())

	tests := []struct {
		name           string
		node           staticNode
		path           string
		expectedRoute  Route
		expectedParams context.Params
	}{
		{
			name:           "Exact Match",
			node:           product,
			path:           "product/item1",
			expectedRoute:  itemRoute,
			expectedParams: append(params, context.Param{Key: "item", Value: "item1"}),
		},
		{
			name:           "Match with SubPath",
			node:           product,
			path:           "product/item1/view",
			expectedRoute:  viewRoute,
			expectedParams: append(params, context.Param{Key: "item", Value: "item1"}),
		},
		{
			name:           "Exact Match with Skip SubPath",
			node:           productSkipSubpath,
			path:           "product/item1/order",
			expectedRoute:  itemRoute,
			expectedParams: append(params, context.Param{Key: "item", Value: "item1"}),
		},
		{
			name:           "No Match with SubPath",
			node:           product,
			path:           "product/item1/order",
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
			if !reflect.DeepEqual(params, tt.expectedParams) {
				t.Errorf("%s: expected params %v, got %v", tt.name, tt.expectedParams, params)
			}
		})
	}
}
