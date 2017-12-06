package gorouter

import "testing"

func TestRootNode(t *testing.T) {
	n := newRoot("")

	equal(t, "", n.id)
	equal(t, "", n.regexpToString())
	equal(t, true, n.isRoot())
	equal(t, true, n.isLeaf())
}

func TestRegexNode(t *testing.T) {
	n := newRoot("")

	regexp := "r([a-z]+)go"
	n.setRegexp(regexp)
	equal(t, regexp, n.regexpToString())

	if !n.regexp.MatchString("rzgo") {
		t.Error("Regexp does not match string")
	}
}

func TestUnknownNodesChild(t *testing.T) {
	n := newRoot("")

	node, params := n.child([]string{"a", "b", "c"})

	if node != nil || params != nil {
		t.Error("Node should return nil values for unknown path")
	}
}
