package executor

import (
	"fmt"
	"github.com/agaffney/crapsh/cmd/builtin"
	"github.com/agaffney/crapsh/core/state"
	"github.com/agaffney/crapsh/parser/ast"
	"github.com/agaffney/crapsh/util"
	"os"
)

type CompleteCommand struct {
	astNode   *ast.CompleteCommand
	Pipelines []*Pipeline
}

func NewCompleteCommand(astNode ast.Node) *CompleteCommand {
	c := &CompleteCommand{astNode: astNode.(*ast.CompleteCommand)}
	if c.astNode != nil {
		for _, pipeline := range c.astNode.Pipelines {
			p := NewPipeline(pipeline)
			c.Pipelines = append(c.Pipelines, p)
		}
	}
	return c
}

func (c *CompleteCommand) Run(state *state.State) error {
	for _, pipeline := range c.Pipelines {
		for _, command := range pipeline.Commands {
			args := command.Words
			//fmt.Printf("Args: %#v\n", args)
			foundBuiltin := false
			for _, b := range builtin.Builtins {
				if b.Name == args[0] {
					foundBuiltin = true
					ret := b.Entrypoint(state, args)
					if false {
						fmt.Printf("returned %d\n", ret)
					}
					break
				}
			}
			if !foundBuiltin {
				paths := util.FindExecutables(args[0], util.SplitPathVar(os.Getenv("PATH")), false)
				fmt.Printf("Paths: %#v\n", paths)
				return fmt.Errorf("%s: command not found\n", args[0])
			}
		}
	}
	return nil
}

type Pipeline struct {
	astNode    *ast.Pipeline
	Commands   []*Command
	AndOrToken int
}

func NewPipeline(astNode ast.Node) *Pipeline {
	p := &Pipeline{astNode: astNode.(*ast.Pipeline)}
	if p.astNode != nil {
		for _, command := range p.astNode.Commands {
			c := NewCommand(command)
			p.Commands = append(p.Commands, c)
		}
		if p.astNode.AndOrToken != nil {
			p.AndOrToken = p.astNode.AndOrToken.Type
		}
	}
	return p
}

type Command struct {
	astNode     *ast.SimpleCommand
	Assignments []*Assignment
	Redirects   []*Redirect
	Words       []string
}

func NewCommand(astNode ast.Node) *Command {
	c := &Command{astNode: astNode.(*ast.SimpleCommand)}
	for _, assignment := range c.astNode.Assignments {
		tmpNode := assignment.(*ast.Assignment)
		a := &Assignment{Name: tmpNode.Var, Value: tmpNode.Value}
		c.Assignments = append(c.Assignments, a)
	}
	for _, redirect := range c.astNode.Redirects {
		tmpNode := redirect.(*ast.IoRedirect)
		r := &Redirect{FileNumber: tmpNode.FileNumber, Operator: tmpNode.Operator, Target: tmpNode.Target}
		c.Redirects = append(c.Redirects, r)
	}
	for _, word := range c.astNode.Words {
		c.Words = append(c.Words, word.GetToken().Value)
	}
	return c
}

type Assignment struct {
	Name  string
	Value string
}

type Redirect struct {
	FileNumber int
	Operator   int
	Target     string
}
