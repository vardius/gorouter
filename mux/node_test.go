package mux

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/vardius/gorouter/v4/context"
	"github.com/vardius/gorouter/v4/middleware"
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

func TestRegexpdNodeMatchRoute(t *testing.T) {
	paramSize := 3
	productRoute := newMockRoute("testproductroute")
	itemRoute := newMockRoute("testitemroute")
	viewRoute := newMockRoute("testviewroute")
	params := make(context.Params, paramSize)

	product := staticNode{name: "product", route: nil, maxParamsSize: uint8(paramSize)}
	product.WithRoute(productRoute)

	// Create a regexp node for "{item}" under "product"
	item := NewNode("{item:item1|item2}", product.MaxParamsSize())
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

	itemSkipSubpath := NewNode("{item:item1|item2}", productSkipSubpath.MaxParamsSize()) // regexpNode
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
			expectedRoute:  item.Route(),
			expectedParams: append(params, context.Param{Key: "item", Value: "item1"}),
		},
		{
			name:           "Match with SubPath",
			node:           product,
			path:           "product/item1/view",
			expectedRoute:  view.Route(),
			expectedParams: append(params, context.Param{Key: "item", Value: "item1"}),
		},
		{
			name:           "Exact Match with Skip SubPath",
			node:           productSkipSubpath,
			path:           "product/item2/order",
			expectedRoute:  item.Route(),
			expectedParams: append(params, context.Param{Key: "item", Value: "item2"}),
		},
		{
			name:           "No Match with SubPath",
			node:           product,
			path:           "product/item1/order",
			expectedRoute:  nil,
			expectedParams: nil,
		},
		{
			name:           "No Match",
			node:           product,
			path:           "product/item3/view",
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

func buildMockMiddlewareFunc(body string) middleware.Middleware {
	fn := func(h middleware.Handler) middleware.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if _, err := w.Write([]byte(body)); err != nil {
				panic(err)
			}
			h.(http.Handler).ServeHTTP(w, r)
		})
	}

	return middleware.WrapperFunc(fn)
}

type middlewareTest struct {
	name           string
	node           Node
	path           string
	expectedResult middleware.Collection
}

func TestStaticNodeMatchMiddleware(t *testing.T) {
	paramSize := 3
	middleware1 := buildMockMiddlewareFunc("1")
	middleware2 := buildMockMiddlewareFunc("2")
	middleware3 := buildMockMiddlewareFunc("3")
	allMiddleware := middleware.NewCollection(middleware1, middleware2, middleware3)

	node1 := NewNode("test", uint8(paramSize))
	node1.PrependMiddleware(allMiddleware)

	node2 := NewNode("test", uint8(paramSize))
	node2.PrependMiddleware(allMiddleware)
	node2.SkipSubPath()

	node3 := NewNode("test", uint8(paramSize))
	node3Middleware := middleware.NewCollection(middleware1, middleware2)
	node3.PrependMiddleware(node3Middleware)
	subpathNode := NewNode("subpath", node3.MaxParamsSize())
	subpathMiddleware := middleware.NewCollection(middleware3)
	subpathNode.PrependMiddleware(subpathMiddleware)
	node3.WithChildren(node3.Tree().withNode(subpathNode).sort())
	node3.WithChildren(node3.Tree().Compile())

	tests := []middlewareTest{
		{
			name:           "StaticNode Exact match",
			node:           node1,
			path:           "test",
			expectedResult: allMiddleware,
		},
		{
			name:           "StaticNode Subpath match with skipSubPath",
			node:           node2,
			path:           "test/subpath",
			expectedResult: allMiddleware,
		},
		{
			name:           "StaticNode Subpath match without skipSubPath",
			node:           node3,
			path:           "test/subpath",
			expectedResult: allMiddleware,
		},
		{
			name:           "StaticNode Match with only prefix",
			node:           node3,
			path:           "testxyz",
			expectedResult: node3Middleware,
		},
		{
			name:           "StaticNode No match",
			node:           node1,
			path:           "nomatch",
			expectedResult: nil,
		},
	}

	runMiddlewareTests(tests, t)
}

func TestWildcardNodeMatchMiddleware(t *testing.T) {
	paramSize := 3
	middleware1 := buildMockMiddlewareFunc("1")
	middleware2 := buildMockMiddlewareFunc("2")
	middleware3 := buildMockMiddlewareFunc("3")
	middleware4 := buildMockMiddlewareFunc("4")
	mw1 := middleware.NewCollection(middleware1)
	mw2 := middleware.NewCollection(middleware2, middleware3)
	mw3 := middleware.NewCollection(middleware4)

	node1 := NewNode("test", uint8(paramSize))
	node1.PrependMiddleware(mw1)
	item := NewNode("{item}", node1.MaxParamsSize()) // wildcardnode
	item.PrependMiddleware(mw2)
	view := NewNode("view", node1.MaxParamsSize()+1)
	view.PrependMiddleware(mw3)
	node1.WithChildren(node1.Tree().withNode(item).sort())
	item.WithChildren(item.Tree().withNode(view).sort())
	node1.WithChildren(node1.Tree().Compile())

	node2 := NewNode("test", uint8(paramSize))
	node2.PrependMiddleware(mw1)
	item2 := NewNode("{item}", node1.MaxParamsSize()) // wildcardnode
	item2.PrependMiddleware(mw2)
	item2.SkipSubPath()
	node2.WithChildren(node2.Tree().withNode(item2).sort())
	node2.WithChildren(node2.Tree().Compile())

	tests := []middlewareTest{
		{
			name:           "WildcardNode Exact match",
			node:           node1,
			path:           "test/item1",
			expectedResult: middleware.NewCollection(middleware1, middleware2, middleware3),
		},
		{
			name:           "WildcardNode Subpath match with skipSubPath",
			node:           node2,
			path:           "test/item2/random",
			expectedResult: middleware.NewCollection(middleware1, middleware2, middleware3),
		},
		{
			name:           "WildcardNode Subpath match without skipSubPath",
			node:           node1,
			path:           "test/item3/view",
			expectedResult: middleware.NewCollection(middleware1, middleware2, middleware3, middleware4),
		},
		{
			name:           "WildcardNode Subpath No match",
			node:           node1,
			path:           "test/item4/nomatch",
			expectedResult: middleware.NewCollection(middleware1, middleware2, middleware3),
		},
	}

	runMiddlewareTests(tests, t)
}

func TestRegexpNodeMatchMiddleware(t *testing.T) {
	paramSize := 3
	middleware1 := buildMockMiddlewareFunc("1")
	middleware2 := buildMockMiddlewareFunc("2")
	middleware3 := buildMockMiddlewareFunc("3")
	middleware4 := buildMockMiddlewareFunc("4")
	mw1 := middleware.NewCollection(middleware1)
	mw2 := middleware.NewCollection(middleware2, middleware3)
	mw3 := middleware.NewCollection(middleware4)

	node1 := NewNode("test", uint8(paramSize))
	node1.PrependMiddleware(mw1)
	item := NewNode("{item:item1|item2}", node1.MaxParamsSize()) // regexpnode
	item.PrependMiddleware(mw2)
	view := NewNode("view", node1.MaxParamsSize()+1)
	view.PrependMiddleware(mw3)
	node1.WithChildren(node1.Tree().withNode(item).sort())
	item.WithChildren(item.Tree().withNode(view).sort())
	node1.WithChildren(node1.Tree().Compile())

	node2 := NewNode("test", uint8(paramSize))
	node2.PrependMiddleware(mw1)
	item2 := NewNode("{item:item1|item2}", node1.MaxParamsSize()) // regexpnode
	item2.PrependMiddleware(mw2)
	item2.SkipSubPath()
	node2.WithChildren(node2.Tree().withNode(item2).sort())
	node2.WithChildren(node2.Tree().Compile())

	tests := []middlewareTest{
		{
			name:           "RegexpNode Exact match",
			node:           node1,
			path:           "test/item1",
			expectedResult: middleware.NewCollection(middleware1, middleware2, middleware3),
		},
		{
			name:           "RegexpNode Subpath match with skipSubPath",
			node:           node2,
			path:           "test/item2/random",
			expectedResult: middleware.NewCollection(middleware1, middleware2, middleware3),
		},
		{
			name:           "RegexpNode Subpath match without skipSubPath",
			node:           node1,
			path:           "test/item1/view",
			expectedResult: middleware.NewCollection(middleware1, middleware2, middleware3, middleware4),
		},
		{
			name:           "RegexpNode Subpath No match",
			node:           node1,
			path:           "test/item2/nomatch",
			expectedResult: middleware.NewCollection(middleware1, middleware2, middleware3),
		},
		{
			name:           "RegexpNode No match",
			node:           node1,
			path:           "test/item3/view",
			expectedResult: middleware.NewCollection(middleware1),
		},
	}

	runMiddlewareTests(tests, t)
}

func runMiddlewareTests(tests []middlewareTest, t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.node.MatchMiddleware(tt.path)
			if len(got) != len(tt.expectedResult) {
				t.Errorf("%s: middleware length mismatch: got= %v, want %v", tt.name, got, tt.expectedResult)
			} else {
				for k, v := range tt.expectedResult {
					// reflect.DeepEqual do not work for function values.
					// hence compare the pointers of functions as a substitute.
					// function pointers are unique to each function, even if the functions have the same code.
					expectedPointer := reflect.ValueOf(v).Pointer()
					gotPointer := reflect.ValueOf(got[k]).Pointer()
					if expectedPointer != gotPointer {
						t.Errorf("%s: middleware mismatch: got= %v, want %v", tt.name, v, got[k])
					}
				}
			}
		})
	}
}
