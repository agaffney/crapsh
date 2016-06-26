package parser

import (
	"fmt"
	"github.com/agaffney/crapsh/lang"
	"github.com/agaffney/crapsh/util"
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
	//	hint           *lang.ParserHint
	//	allowed        []*lang.ParserHint
	entry          *lang.ElementEntry
	position       Position
	element        lang.Element
	parentEndToken *lang.ParserHint
	hintIdx        int
}

func (stack *Stack) Reset() {
	stack.entries = nil
	stack.entries = make([]*StackEntry, 0)
	stack.depth = -1
}

func (stack *Stack) Add(entry *lang.ElementEntry) {
	//	allowed := []*lang.ParserHint{}
	//	if hint.CaptureAll {
	//		// If the current stack entry captures all, we need to use the allowed
	//		// elements from the "parent". We also want to filter ourselves out
	//		for _, foo := range lang.GetElementHints(stack.Cur().hint.AllowedElements) {
	//			if !foo.CaptureAll {
	//				allowed = append(allowed, foo)
	//			}
	//		}
	//		// Use the end token info from the parent
	//		hint.TokenEnd = stack.Cur().hint.TokenEnd
	//		hint.EndOnWhitespace = stack.Cur().hint.EndOnWhitespace
	//	} else {
	//		allowed = lang.GetElementHints(hint.AllowedElements)
	//	}
	//	parentTokenEnd := ""
	//	if stack.depth > MIN_STACK_DEPTH {
	//		if foo := stack.Cur().hint.TokenEnd; foo != "" && !stack.Cur().hint.EndTokenOptional {
	//			parentTokenEnd = foo
	//		} else {
	//			parentTokenEnd = stack.Cur().parentTokenEnd
	//		}
	//	}
	e := &StackEntry{entry: entry, position: stack.parser.Position}
	stack.entries = append(stack.entries, e)
	stack.depth++
	//e.element = stack.parser.newElement()
	fmt.Printf("\n>>> stack[%d] = %#v\n\n", stack.depth, entry)
	//fmt.Printf("  allowed = [\n")
	//for _, foo := range e.allowed {
	//	fmt.Printf("    %#v,\n", foo)
	//}
	//fmt.Printf("  ]\n\n")
}

func (stack *Stack) Remove() {
	//stack.Cur().element.SetContent(buf.String())
	//if stack.depth > MIN_STACK_DEPTH {
	//	stack.Prev().element.AddChild(stack.Cur().element)
	//}
	util.DumpJson(stack.entries[stack.depth].element, "\nremoving element: ")
	stack.entries = stack.entries[:len(stack.entries)-1]
	stack.depth--
	if stack.depth >= MIN_STACK_DEPTH {
		fmt.Printf("\n<<< stack[%d] = %#v\n\n", stack.depth, stack.entries[stack.depth].entry)
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

func (entry *StackEntry) NextHint() *lang.ParserHint {
	if entry.hintIdx < len(entry.entry.ParserData) {
		return entry.entry.ParserData[entry.hintIdx+1]
	}
	return nil
}
