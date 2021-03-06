package ast

import (
	"github.com/agaffney/crapsh/parser/lexer"
)

type Node interface {
	GetName() string
	AddToken(*lexer.Token)
	GetToken() *lexer.Token
	AddChild(Node)
	GetChildren() []Node
}

type NodeBase struct {
	Name  string
	Nodes []Node
}

func NewNode(name string) *NodeBase {
	n := &NodeBase{Name: name}
	return n
}

func (n *NodeBase) GetName() string {
	return n.Name
}

func (n *NodeBase) AddToken(token *lexer.Token) {
	w := NewWord(token)
	n.AddChild(w)
}

func (n *NodeBase) GetToken() *lexer.Token {
	// Not implemented in NodeBase
	return nil
}

func (n *NodeBase) AddChild(node Node) {
	n.Nodes = append(n.Nodes, node)
}

func (n *NodeBase) GetChildren() []Node {
	return n.Nodes
}

type Word struct {
	NodeBase
	Token *lexer.Token
}

func NewWord(token *lexer.Token) Node {
	w := &Word{NodeBase: NodeBase{Name: `Word`}, Token: token}
	return w
}

func (w *Word) GetToken() *lexer.Token {
	return w.Token
}
