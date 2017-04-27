package goserver

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
		t.Error("Regexp does not amtch string")
	}
}
