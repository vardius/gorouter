package mux

import (
	"testing"
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
		t.Fatalf("Expecting: *mux.staticNode. Wrong node type: %T\n", node)
	case *wildcardNode:
		t.Fatalf("Expecting: *mux.staticNode. Wrong node type: %T\n", node)
	}

	node = NewNode("{lang}", 0)

	switch node := node.(type) {
	case *staticNode:
		t.Fatalf("Expecting: *mux.staticNode. Wrong node type: %T\n", node)
	case *regexpNode:
		t.Fatalf("Expecting: *mux.staticNode. Wrong node type: %T\n", node)
	}

}
