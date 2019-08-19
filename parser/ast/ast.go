package ast

import (
	"github.com/agaffney/crapsh/parser/lexer"
)

type Node interface {
	GetName() string
	GetParent() Node
	GetTokens() []*lexer.Token
	AddToken(*lexer.Token)
	AddChild(Node)
}

type NodeBase struct {
	Name   string
	Parent Node
	Tokens []*lexer.Token
	Nodes  []Node
}

func NewNode() *NodeBase {
	n := &NodeBase{}
	n.Tokens = make([]*lexer.Token, 0)
	n.Nodes = make([]Node, 0)
	return n
}

func (n *NodeBase) GetName() string {
	return n.Name
}

func (n *NodeBase) GetParent() Node {
	return n.Parent
}

func (n *NodeBase) GetTokens() []*lexer.Token {
	return n.Tokens
}

func (n *NodeBase) AddToken(token *lexer.Token) {
	n.Tokens = append(n.Tokens, token)
}

func (n *NodeBase) AddChild(node Node) {
	n.Nodes = append(n.Nodes, node)
}
