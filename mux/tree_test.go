package mux

import "testing"

func TestInsertNilNode(t *testing.T) {
	tree := NewTree()
	tree.Insert(nil)

	if tree.idsLen > 0 || len(tree.regexps) > 0 || tree.wildcard != nil {
		t.Error("Tree should not insert nil node")
	}
}

func TestTreeInsertWildcardNodePanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Tree should panic on wildcard insert if one already exist")
		}
	}()

	n1 := NewRoot("")
	n1.isWildcard = true

	n2 := NewRoot("")
	n2.isWildcard = true

	tree := NewTree()
	tree.Insert(n1)
	tree.Insert(n2)
}

func TestTreeGetRegexNodeById(t *testing.T) {
	n := NewRoot("")
	n.SetRegexp("r([a-z]+)go")

	tree := NewTree()
	tree.Insert(n)

	c := tree.GetByID("rego")

	if c != n {
		t.Error("Tree should match regex node by ID")
	}
}

func TestGetTreeNodeByEmptyPath(t *testing.T) {
	tree := NewTree()
	n, _, _ := tree.GetByPath("")

	if n != nil {
		t.Error("Tree should return nil node for empty path")
	}
}
