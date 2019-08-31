package ast

import (
	"github.com/agaffney/crapsh/parser/lexer"
	"github.com/agaffney/crapsh/parser/tokens"
	"strconv"
	"strings"
)

type CompleteCommand struct {
	NodeBase
	Background bool // whether to background the command
	Pipelines  []Node
}

func NewCompleteCommand() Node {
	c := &CompleteCommand{NodeBase: NodeBase{Name: `CompleteCommand`}}
	c.Pipelines = make([]Node, 0)
	return c
}

func (c *CompleteCommand) AddChild(node Node) {
	switch node.GetName() {
	case `and_or`:
		for _, tmpNode := range node.GetChildren() {
			c.AddChild(tmpNode)
		}
	case `Pipeline`:
		c.Pipelines = append(c.Pipelines, node)
	case `separator_op`:
		for _, tmpNode := range node.GetChildren() {
			c.AddChild(tmpNode)
		}
	case `Word`:
		if node.GetToken().Type == tokens.TOKEN_AND {
			c.Background = true
		} else {
			c.Pipelines[len(c.Pipelines)-1].AddChild(node)
		}
	}
}

func (c *CompleteCommand) GetChildren() []Node {
	return c.Pipelines
}

type Pipeline struct {
	NodeBase
	Bang       bool // whether ! appears before pipeline
	Commands   []Node
	AndOrToken *lexer.Token // token with &&/|| suffix
}

func NewPipeline() Node {
	p := &Pipeline{NodeBase: NodeBase{Name: `Pipeline`}}
	p.Commands = make([]Node, 0)
	return p
}

func (p *Pipeline) AddToken(token *lexer.Token) {
	// Empty function to discard pipe character
}

func (p *Pipeline) AddChild(node Node) {
	switch node.GetName() {
	case `Word`:
		if node.GetToken().Type == tokens.TOKEN_BANG {
			p.Bang = true
		} else {
			p.AndOrToken = node.GetToken()
		}
	case `command`:
		for _, tmpNode := range node.GetChildren() {
			p.AddChild(tmpNode)
		}
	case `linebreak`:
		// Purposely empty
	default:
		p.Commands = append(p.Commands, node)
	}
}

func (p *Pipeline) GetChildren() []Node {
	return p.Commands
}

type Assignment struct {
	NodeBase
	Var   string
	Value string
}

func NewAssignment() Node {
	a := &Assignment{NodeBase: NodeBase{Name: `Assignment`}}
	return a
}

func (a *Assignment) AddChild(node Node) {
	switch node.GetName() {
	case `cmd_prefix`:
		for _, tmpNode := range node.GetChildren() {
			a.AddChild(tmpNode)
		}
	default:
		token := node.GetToken()
		parts := strings.SplitN(token.Value, `=`, 2)
		a.Var = parts[0]
		a.Value = parts[1]
	}
}

type IoRedirect struct {
	NodeBase
	FileNumber int
	Operator   int
	Target     string
}

func NewIoRedirect() Node {
	i := &IoRedirect{NodeBase: NodeBase{Name: `IoRedirect`}, FileNumber: -1}
	return i
}

func (i *IoRedirect) AddChild(node Node) {
	switch node.GetName() {
	case `io_file`:
		for _, tmpNode := range node.GetChildren() {
			i.AddChild(tmpNode)
		}
	case `io_here`:
		for _, tmpNode := range node.GetChildren() {
			i.AddChild(tmpNode)
		}
	case `Word`:
		token := node.GetToken()
		switch token.Type {
		case tokens.TOKEN_IO_NUMBER:
			num, _ := strconv.Atoi(token.Value)
			i.FileNumber = num
		case tokens.TOKEN_DLESS, tokens.TOKEN_LESSAND, tokens.TOKEN_DLESSDASH, tokens.TOKEN_LESS:
			if i.FileNumber == -1 {
				// STDIN
				i.FileNumber = 0
			}
			i.Operator = token.Type
		case tokens.TOKEN_DGREAT, tokens.TOKEN_GREATAND, tokens.TOKEN_CLOBBER, tokens.TOKEN_GREAT:
			if i.FileNumber == -1 {
				// STDOUT
				i.FileNumber = 1
			}
			i.Operator = token.Type
			//TOKEN_LESSGREAT // <>
		default:
			i.Target = token.Value
		}
	}
}

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
	switch node.GetName() {
	case `IoRedirect`:
		c.Redirects = append(c.Redirects, node)
	case `Word`:
		token := node.GetToken()
		switch token.Type {
		case tokens.TOKEN_ASSIGNMENT_WORD:
			a := NewAssignment()
			a.AddChild(node)
			c.Assignments = append(c.Assignments, a)
		default:
			c.Words = append(c.Words, node)
		}
	case `cmd_prefix`, `cmd_suffix`:
		for _, tmpNode := range node.GetChildren() {
			c.AddChild(tmpNode)
		}
	default:
		c.Nodes = append(c.Nodes, node)
	}
}

func (c *SimpleCommand) GetChildren() []Node {
	var ret []Node
	ret = append(ret, c.Assignments...)
	ret = append(ret, c.Words...)
	ret = append(ret, c.Redirects...)
	return ret
}
