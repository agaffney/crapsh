package ast

import (
	"fmt"
	"github.com/agaffney/crapsh/parser/lexer"
)

type Command struct {
	NodeBase
}

func NewCommand() Node {
	c := &Command{NodeBase: NodeBase{Name: `Command`}}
	return c
}

type IoRedirect struct {
	NodeBase
	FileNumber int
	Operator   int
}

/*
func (i *IoRedirect) AddChild(node Node) {
	switch node.GetName() {
	case `io_redirect`:
		for _, tmpNode := range node.GetChildren() {
			i.AddChild(tmpNode)
		}
	case `word`:
		if node.
	}
}
*/

type SimpleCommand struct {
	NodeBase
	Assignments []Node
	Redirects   []Node
	Words       []Node
}

func NewSimpleCommand() Node {
	c := &SimpleCommand{NodeBase: NodeBase{Name: `SimpleCommand`}}
	c.Assignments = make([]Node, 0)
	c.Redirects = make([]Node, 0)
	c.Words = make([]Node, 0)
	return c
}

func (c *SimpleCommand) AddToken(token *lexer.Token) {
	w := NewWord(token)
	c.AddChild(w)
}

func (c *SimpleCommand) AddChild(node Node) {
	fmt.Printf("SimpleCommand.AddChild(): node = %#v\n", node)
	switch node.GetName() {
	case `io_redirect`:
		for _, tmpNode := range node.GetChildren() {
			c.AddChild(tmpNode)
		}
	case `io_file`:
		c.Redirects = append(c.Redirects, node)
	case `io_here`:
		c.Redirects = append(c.Redirects, node)
	case `word`:
		c.Words = append(c.Words, node)
	case `cmd_suffix`:
		for _, tmpNode := range node.GetChildren() {
			c.AddChild(tmpNode)
		}
	default:
		c.Nodes = append(c.Nodes, node)
	}

}
