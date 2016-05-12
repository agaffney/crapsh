package tree

import (
	"github.com/agaffney/crapsh/parser/tokens"
)

type StringNode struct {
	*BaseNode
}

func init() {
	add_node_type(NodeType{
		Name:      "String",
		Container: true,
		Token:     tokens.DOUBLE_QUOTE,
		Factory:   NewString,
	})
}

func NewString(parent Node) Node {
	return &StringNode{BaseNode: NewBaseNode(parent)}
}
