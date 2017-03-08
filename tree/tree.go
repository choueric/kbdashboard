package tree

import (
	"fmt"
	"io"
)

type Node struct {
	Value    interface{}
	subNodes []*Node
}

// New creates and returns a root node of one tree with val.
func New(val interface{}) *Node {
	return &Node{Value: val}
}

// AddSubNode adds a sub node to n and return the sub node.
func (n *Node) AddSubNode(val interface{}) *Node {
	sub := &Node{Value: val}
	n.subNodes = append(n.subNodes, sub)
	return sub
}

// Walk goes trough the tree by f, first at n and then its sub nodes.
func (n *Node) Walk(f func(*Node)) {
	f(n)
	for _, sub := range n.subNodes {
		sub.Walk(f)
	}
}

func (n *Node) doPrintTree(w io.Writer, f func(*Node) string, depth int, hasSibling []bool) {
	for i := 0; i < depth; i++ {
		if i != depth-1 {
			if hasSibling[i] {
				fmt.Fprintf(w, "│   ")
			} else {
				fmt.Fprintf(w, "    ")
			}
		} else {
			if hasSibling[i] {
				fmt.Fprintln(w, "├── "+f(n))
			} else {
				fmt.Fprintln(w, "└── "+f(n))
			}
		}
	}

	if depth == 0 {
		fmt.Fprintln(w, f(n))
	}

	length := len(n.subNodes)
	for i, subNode := range n.subNodes {
		if i == length-1 {
			hasSibling = append(hasSibling, false)
		} else {
			hasSibling = append(hasSibling, true)
		}

		subNode.doPrintTree(w, f, depth+1, hasSibling)
		hasSibling = append(hasSibling[:len(hasSibling)-1])
	}
}

// PrintTree prints the tree graphic started from node n.
func (n *Node) PrintTree(w io.Writer, f func(*Node) string) {
	var hasSibling []bool
	n.doPrintTree(w, f, 0, hasSibling)
}
