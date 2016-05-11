package tree

type Node interface {
	Add_child(Node)
}

type BaseNode struct {
	parent   *Node
	children []Node
	Line     uint
	Offset   uint
}

type TopNode struct {
	*BaseNode
}

type GenericNode struct {
	*BaseNode
	Content string
}

func NewNode(parent Node) *BaseNode {
	return &BaseNode{parent: &parent}
}

func (n *BaseNode) Add_child(child Node) {
	if n.children == nil {
		n.children = make([]Node, 0)
	}
	n.children = append(n.children, child)
}
