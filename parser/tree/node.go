package tree

type Node interface {
	Add_child(Node)
	Set_content(string)
}

type FactoryFunc func(Node) Node

type NodeType struct {
	Name      string
	Container bool
	Token     string
	TokenEnd  string
	Factory   FactoryFunc
}

var Node_types []NodeType

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

func init() {
	Node_types = make([]NodeType, 0)
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

func add_node_type(nt NodeType) {
	Node_types = append(Node_types, nt)
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
