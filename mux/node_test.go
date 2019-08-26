package mux

import (
	"reflect"
	"testing"
)

func equal(t *testing.T, expected, actual interface{}) bool {
	if !areEqual(expected, actual) {
		t.Errorf("Asserts are not equal. Expected: %v, Actual: %v", expected, actual)

		return false
	}

	return true
}

func notEqual(t *testing.T, expected, actual interface{}) bool {
	if areEqual(expected, actual) {
		t.Errorf("Asserts are equal. Expected: %v, Actual: %v", expected, actual)

		return false
	}

	return true
}

func areEqual(expected, actual interface{}) bool {
	if expected == nil {
		return isNil(actual)
	}

	if actual == nil {
		return isNil(expected)
	}

	return reflect.DeepEqual(expected, actual)
}

func isNil(value interface{}) bool {
	defer func() { recover() }()
	return value == nil || reflect.ValueOf(value).IsNil()
}

func TestRootNode(t *testing.T) {
	n := NewRoot("")

	equal(t, "", n.id)
	equal(t, "", n.regexpToString())
	equal(t, true, n.IsRoot())
	equal(t, true, n.IsLeaf())
}

func TestRegexNode(t *testing.T) {
	n := NewRoot("")

	regexp := "r([a-z]+)go"
	n.SetRegexp(regexp)
	equal(t, regexp, n.regexpToString())

	if !n.regexp.MatchString("rzgo") {
		t.Error("Regexp does not match string")
	}
}

func TestUnknownNodesChild(t *testing.T) {
	n := NewRoot("")

	node, params := n.GetChild([]string{"a", "b", "c"})

	if node != nil || params != nil {
		t.Error("Node should return nil values for unknown path")
	}
}
