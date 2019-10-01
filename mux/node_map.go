package mux

import (
	"sort"
)

type NodeMap struct {
	len   int
	ids   []string
	nodes []*Node
}

// NewNodeMap provides new node map
func NewNodeMap() *NodeMap {
	return &NodeMap{
		ids:   make([]string, 0),
		nodes: make([]*Node, 0),
	}
}

func (m *NodeMap) Append(n *Node) {
	m.ids = append(m.ids, n.id)
	m.nodes = append(m.nodes, n)
	m.len++

	sort.Slice(m.ids, func(i, j int) bool {
		return len(m.ids[i]) < len(m.ids[j])
	})

	sort.Slice(m.nodes, func(i, j int) bool {
		return len(m.nodes[i].id) < len(m.nodes[j].id)
	})
}

func (m *NodeMap) Find(id string) *Node {
	if m.len == 0 {
		return nil
	}

	// there's an overhead to looking up values in a map, whereas as slice is directly indexing the values
	for index, internalID := range m.ids {
		if id == internalID {
			return m.nodes[index]
		}
	}

	return nil
}
