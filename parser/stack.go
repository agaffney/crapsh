package parser

import (
	//"fmt"
	//"github.com/agaffney/crapsh/util"
	"github.com/agaffney/crapsh/parser/ast"
	"github.com/agaffney/crapsh/parser/rules/grammar"
)

const (
	MIN_STACK_DEPTH = 0
)

type Stack struct {
	entries []*StackEntry
	depth   int
	parser  *Parser
}

type StackEntry struct {
	rule                  *grammar.GrammarRule
	astNode               ast.Node
	hintIdx               int
	final                 bool
	allowNextWordReserved bool
}

func (stack *Stack) Reset() {
	stack.entries = make([]*StackEntry, 0)
	stack.depth = -1
}

func (stack *Stack) Add(rule *grammar.GrammarRule) {
	e := &StackEntry{rule: rule}
	stack.entries = append(stack.entries, e)
	stack.depth++
	//fmt.Printf("\n>>> stack[%d] = %#v\n\n", stack.depth, rule)
}

func (stack *Stack) Remove() {
	//stack.Cur().element.SetContent(buf.String())
	//if stack.depth > MIN_STACK_DEPTH {
	//	stack.Prev().element.AddChild(stack.Cur().element)
	//}
	//util.DumpJson(stack.entries[stack.depth].element, "\nremoving element: ")
	stack.entries = stack.entries[:len(stack.entries)-1]
	stack.depth--
	/*
		if stack.depth >= MIN_STACK_DEPTH {
			fmt.Printf("\n<<< stack[%d] = %#v\n\n", stack.depth, stack.entries[stack.depth].rule)
		}
	*/
}

func (stack *Stack) Cur() *StackEntry {
	if stack.depth >= MIN_STACK_DEPTH {
		return stack.entries[stack.depth]
	}
	return nil
}

func (stack *Stack) Prev() *StackEntry {
	if stack.depth-1 >= MIN_STACK_DEPTH {
		return stack.entries[stack.depth-1]
	}
	return nil
}
