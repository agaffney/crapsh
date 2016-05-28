package parser

import (
	"github.com/agaffney/crapsh/lang"
)

type Stack struct {
	entries []*StackEntry
	depth   int
	parser  *Parser
}

type StackEntry struct {
	hint           *lang.ParserHint
	allowed        []*lang.ParserHint
	position       Position
	element        lang.Element
	parentTokenEnd string
}

func (stack *Stack) Reset() {
	stack.entries = nil
}

func (stack *Stack) Add(hint *lang.ParserHint) {
	if stack.entries == nil {
		stack.entries = make([]*StackEntry, 0)
		stack.depth = -1
	}
	allowed := []*lang.ParserHint{}
	if hint.CaptureAll {
		// If the current stack entry captures all, we need to use the allowed
		// elements from the "parent". We also want to filter ourselves out
		for _, foo := range lang.GetElementHints(stack.Cur().hint.AllowedElements) {
			if !foo.CaptureAll {
				allowed = append(allowed, foo)
			}
		}
		// Use the end token info from the parent
		hint.TokenEnd = stack.Cur().hint.TokenEnd
		hint.EndOnWhitespace = stack.Cur().hint.EndOnWhitespace
	} else {
		allowed = lang.GetElementHints(hint.AllowedElements)
	}
	parentTokenEnd := ""
	if stack.depth > MIN_STACK_DEPTH {
		if foo := stack.Cur().hint.TokenEnd; foo != "" && !stack.Cur().hint.EndTokenOptional {
			parentTokenEnd = foo
		} else {
			parentTokenEnd = stack.Cur().parentTokenEnd
		}
	}
	e := &StackEntry{hint, allowed, stack.parser.Position, nil, parentTokenEnd}
	stack.entries = append(stack.entries, e)
	stack.depth++
	e.element = stack.parser.newElement()
	//fmt.Printf(">>> stack[%d] = %#v\n\n", stack.depth, hint)
	//fmt.Printf("  allowed = [\n")
	//for _, foo := range e.allowed {
	//	fmt.Printf("    %#v,\n", foo)
	//}
	//fmt.Printf("  ]\n\n")
}

func (stack *Stack) Remove(buf *Buffer) {
	stack.Cur().element.SetContent(buf.String())
	if stack.depth > MIN_STACK_DEPTH {
		stack.Prev().element.AddChild(stack.Cur().element)
	}
	stack.entries = stack.entries[:len(stack.entries)-1]
	stack.depth--
	if stack.depth >= MIN_STACK_DEPTH {
		//fmt.Printf("\nstack[%d] = %#v\n\n", stack.depth, p.stack[stack.depth].hint)
	}
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
