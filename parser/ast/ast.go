package ast

import (
	"github.com/agaffney/crapsh/parser/lexer"
)

type Node interface {
	GetName() string
	//GetParent() Node
	AddToken(*lexer.Token)
	AddChild(Node)
}

type NodeBase struct {
	Name string
	//Parent Node
	Nodes []Node
}

func NewNode(name string) *NodeBase {
	n := &NodeBase{Name: name}
	n.Nodes = make([]Node, 0)
	return n
}

func (n *NodeBase) GetName() string {
	return n.Name
}

/*
func (n *NodeBase) GetParent() Node {
	return n.Parent
}
*/

func (n *NodeBase) AddToken(token *lexer.Token) {
	w := NewWord(token)
	n.AddChild(w)
}

func (n *NodeBase) AddChild(node Node) {
	n.Nodes = append(n.Nodes, node)
}

type Word struct {
	NodeBase
	Token *lexer.Token
}

func NewWord(token *lexer.Token) Node {
	w := &Word{NodeBase: NodeBase{Name: `word`}, Token: token}
	return w
}
