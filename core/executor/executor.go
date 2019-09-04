package executor

import (
	"github.com/agaffney/crapsh/parser/ast"
)

type Executor struct {
}

func New() *Executor {
	e := &Executor{}
	return e
}

func (e *Executor) CommandFromAst(astNode ast.Node) *CompleteCommand {
	c := NewCompleteCommand(astNode)
	return c
}
