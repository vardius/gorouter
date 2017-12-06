package gorouter

import "testing"

func TestInsertNilNode(t *testing.T) {
	t := newTree()
	t.insert(nil)

	if t.idsLen > 0 || len(t.regexps) > 0 || t.wildcard != nil {
		t.Error("Tree should not insert nil node")
	}
}

func TestTreeInsertWildcardNodePanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Tree should panic on wildcard insert if one already exist")
		}
	}()

	n1 := newRoot("")
	n1.isWildcard = true

	n2 := newRoot("")
	n2.isWildcard = true

	t := newTree()
	t.insert(n1)
	t.insert(n2)
}

func TestTreeGetRegexNodeById(t *testing.T) {
	n := newRoot("")
	n.setRegexp("r([a-z]+)go")

	t := newTree()
	t.insert(n)

	c := t.byID("rego")

	if c == nil {
		t.Error("Tree should match regex node by ID")
	}
}

func TestGetTreeNodeByEmptyPath(t *testing.T) {
	t := newTree()
	n := t.byPath("")

	if n != nil {
		t.Error("Tree should return nil node for empty path")
	}
}

func TestTreeMerge(t *testing.T) {
	t1 := newTree()
	t2 := newTree()

	n1 := newRoot("")
	n1.isWildcard = true
	t2.insert(n1)

	n2 := newRoot("")
	n2.setRegexp("r([a-z]+)go")
	t2.insert(n2)

	n3 := newRoot("")
	t2.insert(n3)

	if t.idsLen == 0 || len(t.regexps) == 0 || t.wildcard == nil {
		t.Error("Tree should merge sub tree correctly")
	}
}
