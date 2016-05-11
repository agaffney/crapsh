package tree

type Node interface {
	Add_child(Node)
	Set_content(string)
}

type BaseNode struct {
	parent   *Node
	children []Node
	Line     uint
	Offset   uint
	Content  string
}

type TopNode struct {
	*BaseNode
}

type GenericNode struct {
	*BaseNode
}

func NewBaseNode(parent Node) *BaseNode {
	return &BaseNode{parent: &parent}
}

func NewTop() Node {
	return &TopNode{BaseNode: NewBaseNode(nil)}
}

func NewGeneric(parent Node) Node {
	return &GenericNode{BaseNode: NewBaseNode(parent)}
}

func (n *BaseNode) Add_child(child Node) {
	if n.children == nil {
		n.children = make([]Node, 0)
	}
	n.children = append(n.children, child)
}

func (n *BaseNode) Set_content(content string) {
	n.Content = content
}
