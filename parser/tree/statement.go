package tree

type StatementNode struct {
	*BaseNode
}

func init() {
	add_node_type(NodeType{
		Name:    "Statement",
		Factory: NewStatement,
	})
}

func NewStatement(parent Node) Node {
	return &StatementNode{BaseNode: NewBaseNode(parent)}
}
