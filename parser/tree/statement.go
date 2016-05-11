package tree

type StatementNode struct {
	*BaseNode
}

func NewStatement(parent Node) Node {
	return &StatementNode{BaseNode: NewBaseNode(parent)}
}
