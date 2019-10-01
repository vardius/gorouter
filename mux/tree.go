package mux

// NewTree provides new node tree
func NewTree() *Tree {
	return &Tree{
		statics: NewNodeMap(),
		regexps: make([]*Node, 0),
	}
}

// Tree of routes
type Tree struct {
	wildcard *Node
	statics  *NodeMap
	regexps  []*Node
}

func (t *Tree) StaticNodes() []*Node {
	return t.statics.nodes
}

func (t *Tree) RegexpNodes() []*Node {
	return t.regexps
}

func (t *Tree) WildcardNode() *Node {
	return t.wildcard
}

func (t *Tree) Insert(node *Node) {
	if node == nil {
		return
	}

	if t.wildcard != nil {
		panic("Tree already contains a wildcard node!")
	}

	if node.isRegexp {
		t.regexps = append(t.regexps, node)
	} else if node.isWildcard {
		t.wildcard = node

		// wildcard node will match every case, reset other collections
		t.statics = NewNodeMap()
		t.regexps = make([]*Node, 0)
	} else {
		t.statics.Append(node)
	}
}

// GetByID gets node by ID
// this method is used when inserting new nodes
func (t *Tree) GetByID(id string) *Node {
	if id == "" {
		return nil
	}

	if t.statics.len > 0 {
		for _, staticNode := range t.statics.nodes {
			if staticNode.id == id {
				return staticNode
			}
		}
	}

	if len(t.regexps) > 0 {
		for _, regexpNode := range t.regexps {
			if regexpNode.id == id {
				return regexpNode
			}
		}
	}

	if t.wildcard != nil && t.wildcard.id == id {
		return t.wildcard
	}

	return nil
}

// Find finds node by path part
func (t *Tree) Find(pathPart string) *Node {
	if pathPart == "" {
		return nil
	}

	if t.wildcard != nil {
		return t.wildcard
	}

	staticNode := t.statics.Find(pathPart)
	if staticNode != nil {
		return staticNode
	}

	for _, regexpNode := range t.regexps {
		if regexpNode.regexp.MatchString(pathPart) {
			return regexpNode
		}
	}

	return nil
}
