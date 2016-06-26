package parser

import (
	"bufio"
	"fmt"
	"github.com/agaffney/crapsh/lang"
	"github.com/agaffney/crapsh/util"
	"io"
)

type Parser struct {
	Position
	input    *bufio.Reader
	stack    *Stack
	buf      *Buffer
	tokenBuf []*Token
	tokenIdx int
	LineChan chan lang.Element
	Error    error
}

type Position struct {
	Line       uint
	Offset     uint
	LineOffset uint
}

func NewParser() *Parser {
	parser := &Parser{}
	parser.LineChan = make(chan lang.Element)
	parser.stack = &Stack{parser: parser}
	parser.buf = NewBuffer(nil)
	return parser
}

func (p *Parser) Parse(input io.Reader) {
	r := bufio.NewReader(input)
	p.input = r
	p.Line = 1
	p.LineOffset = 0
	p.Offset = 0
	p.tokenBuf = make([]*Token, 0)
	p.tokenIdx = -1
	go func() {
		for {
			// Reset the hint stack
			p.stack.Reset()
			//line, err := p.scanTokens()
			ok, err := p.parseElement("Line")
			//util.DumpJson(line, "Line: ")
			if err != nil {
				p.Error = err
				break
			}
			// EOF
			if !ok {
				break
			}
			//p.LineChan <- line
		}
		close(p.LineChan)
	}()
}

func (p *Parser) parseHandleHint(hint *lang.ParserHint) (bool, error) {
	//util.DumpObject(hint, "parseHandleHint(): hint = ")
	var err error
	var count int
	var ok bool
	var origTokenIdx int
	for {
		ok = false
		origTokenIdx = p.getTokenIdx()
		switch {
		case hint.Type == lang.HINT_TYPE_ELEMENT:
			ok, err = p.parseElement(hint.Name)
		case hint.Type == lang.HINT_TYPE_ANY:
			ok, err = p.parseAny(hint.Members)
		case hint.Type == lang.HINT_TYPE_GROUP:
			ok, err = p.parseGroup(hint.Members, false)
		case hint.Type == lang.HINT_TYPE_TOKEN:
			ok, err = p.parseToken(hint)
		default:
			return ok, fmt.Errorf("Unhandled hint type: %d\n", hint.Type)
		}
		//util.DumpObject(ok, "parseHandleHint(): ok = ")
		if err != nil {
			return ok, err
		}
		if !ok {
			p.setTokenIdx(origTokenIdx)
			if hint.Optional {
				return true, nil
			}
			if hint.Many && count > 0 {
				return true, nil
			}
			return false, nil
		}
		if !hint.Many {
			break
		}
		count++
	}
	return true, nil
}

func (p *Parser) parseToken(hint *lang.ParserHint) (bool, error) {
	if err := p.nextToken(); err != nil {
		if err == io.EOF {
			return false, nil
		}
		return false, err
	}
	util.DumpObject(hint, "parseToken(): hint = ")
	util.DumpObject(p.curToken(), "parseToken(): curToken = ")
	foo := p.curToken()
	if hint.Name == `` || hint.Name == foo.Type {
		e := lang.NewGeneric(`Token`)
		e.Content = foo.Value
		p.stack.Cur().element.AddChild(e)
		return true, nil
	}
	return false, nil
}

func (p *Parser) parseGroup(hints []*lang.ParserHint, updateHintIdx bool) (bool, error) {
	for idx, hint := range hints {
		//util.DumpObject(hint, "parseGroup(): hint = ")
		if updateHintIdx {
			p.stack.Cur().hintIdx = idx
		}
		ok, err := p.parseHandleHint(hint)
		if err != nil {
			return ok, err
		}
		if !ok {
			return ok, nil
		}
	}
	return true, nil
}

func (p *Parser) parseElement(element string) (bool, error) {
	var parentEndToken *lang.ParserHint
	entry := lang.GetElementEntry(element)
	if entry == nil {
		return false, nil
	}
	// Check next hint for current stack entry to set as the parentEndToken
	// on the new stack entry
	if p.stack.Cur() != nil {
		if nextHint := p.stack.Cur().NextHint(); nextHint != nil {
			if nextHint.Type == lang.HINT_TYPE_TOKEN {
				parentEndToken = nextHint
			}
		}
	}
	//util.DumpObject(entry, "parseElement(): entry = ")
	p.stack.Add(entry)
	e := lang.NewGeneric(entry.Name)
	//fmt.Printf("%#v\n", e)
	if p.stack.Cur().entry.Factory != nil {
		foo := p.stack.Cur().entry.Factory(e)
		//fmt.Printf("%s\n", foo)
		p.stack.Cur().element = foo
	}
	p.stack.Cur().parentEndToken = parentEndToken
	p.stack.Cur().element = e
	ok, err := p.parseGroup(entry.ParserData, true)
	if err != nil {
		return false, err
	}
	p.stack.Remove()
	if ok && p.stack.Cur() != nil {
		p.stack.Cur().element.AddChild(e)
	}
	return ok, nil
}

func (p *Parser) parseAny(hints []*lang.ParserHint) (bool, error) {
	for _, hint := range hints {
		//util.DumpObject(hint, "parseAny(): hint = ")
		ok, err := p.parseHandleHint(hint)
		if err != nil {
			return ok, err
		}
		if ok {
			return ok, nil
		}
	}
	return false, nil
}

//func (p *Parser) GetNextLine() (lang.Element, error) {
//	escape := false
//	// Reset the hint stack
//	p.stack.Reset()
//	p.stack.Add(lang.GetElementHints([]string{"Line"})[0])
//	line_element := p.stack.Cur().element
//	for {
//		c, err := p.nextRune()
//		if err != nil {
//			if err != io.EOF {
//				return nil, err
//			}
//			break
//		}
//		//fmt.Printf(">>> Stack item (%d): %#v\n\t%#v\n", stack.depth+1, stack.Cur(), stack.Cur().hint)
//		//fmt.Printf("Line %d, offset %d, overall offset %d: %#U\n", p.Line, p.LineOffset, p.Offset, c)
//		if c == '\\' && !p.stack.Cur().hint.IgnoreEscapes {
//			escape = !escape
//			p.buf.WriteRune(c)
//			// Explicitly skip to the next iteration so we don't hit
//			// the code below to turn off the 'escape' flag
//			continue
//		}
//		if c == '\n' {
//			p.nextLine()
//			if escape {
//				p.buf.WriteRune(c)
//			} else {
//				if p.stack.Cur().hint.TokenEnd == "" || p.stack.Cur().hint.EndTokenOptional {
//					p.stack.Remove(p.buf)
//					//fmt.Println("removing from stack due to newline")
//					p.buf.Reset()
//					break
//				} else {
//					p.buf.WriteRune(c)
//				}
//			}
//			continue
//		}
//		if unicode.IsSpace(c) && !escape && p.stack.Cur().hint.EndOnWhitespace {
//			if p.stack.Cur().hint.CaptureAll {
//				// We're using the EndOnWhitespace value from our "parent", so if it's found,
//				// we should remove the CaptureAll element from the stack
//				p.stack.Remove(p.buf)
//				//fmt.Println("removing from stack due to CaptureAll and whitespace")
//				p.buf.Reset()
//			}
//			p.stack.Remove(p.buf)
//			//fmt.Println("removing from stack due to whitespace")
//			p.buf.Reset()
//			continue
//		}
//		// Add new character to buf for checking start/end tokens
//		p.buf.WriteRune(c)
//		//fmt.Printf("buf = %#v\n", buf.String())
//		if escape == false {
//			if p.buf.checkForToken(p.stack.Cur().hint.TokenEnd) {
//				// Remove start token from buf
//				p.buf.removeBytes(len(p.stack.Cur().hint.TokenEnd))
//				if p.stack.Cur().hint.CaptureAll {
//					// We're using the end token from our "parent", so if it's found,
//					// we should remove the CaptureAll element from the stack
//					p.stack.Remove(p.buf)
//					//fmt.Println("removing from stack due to CaptureAll and finding end token")
//					p.buf.Reset()
//				}
//				p.stack.Remove(p.buf)
//				//fmt.Println("removing from stack due to finding end token")
//				p.buf.Reset()
//				continue
//			}
//			parentTokenEnd := p.stack.Cur().parentTokenEnd
//			if p.buf.checkForToken(parentTokenEnd) {
//				// Remove end token from buf
//				p.buf.removeBytes(len(p.stack.Cur().parentTokenEnd))
//				if p.stack.Cur().hint.CaptureAll {
//					// We're using the end token from our "parent", so if it's found,
//					// we should remove the CaptureAll element from the stack
//					p.stack.Remove(p.buf)
//					//fmt.Println("removing from stack due to CaptureAll and finding end token")
//					p.buf.Reset()
//				}
//				// Remove items from stack that don't have a required end token
//				for {
//					if p.stack.Cur().hint.EndTokenOptional || (p.stack.Cur().parentTokenEnd != "" && p.stack.Cur().hint.TokenEnd == "" && p.stack.Cur().parentTokenEnd == parentTokenEnd) {
//						p.stack.Remove(p.buf)
//						//fmt.Println("removing extra item from stack due to finding end token")
//						p.buf.Reset()
//					} else {
//						break
//					}
//				}
//				// Put last character of token back in buffer for re-discovery
//				p.unreadRune()
//				p.buf.removeBytes(utf8.RuneLen(c))
//				continue
//			}
//		}
//		found := false
//		for _, hint := range p.stack.Cur().allowed {
//			if hint.SkipCapture || p.buf.checkForToken(hint.TokenStart) || hint.CaptureAll {
//				// Remove start token from buf
//				p.buf.removeBytes(len(hint.TokenStart))
//				if p.stack.Cur().hint.CaptureAll {
//					// We're using the allowed elements from our "parent", so if one is found,
//					// we should remove the CaptureAll element from the stack
//					p.stack.Remove(p.buf)
//					//fmt.Println("removing from stack due to CaptureAll and finding start token")
//					p.buf.Reset()
//				}
//				p.stack.Add(hint)
//				if hint.SkipCapture {
//					p.unreadRune()
//					// Remove last rune from buffer
//					p.buf.removeBytes(utf8.RuneLen(c))
//				}
//				found = true
//				break
//			}
//		}
//		if !found && !p.stack.Cur().hint.CaptureAll {
//			return nil, fmt.Errorf("line %d, pos %d: unexpected character `%c'", p.Position.Line, p.Position.LineOffset, c)
//		}
//		// Reset the 'escape' flag
//		escape = false
//	}
//	// Remove any stack items that allow ending on EOF
//	for p.stack.depth >= MIN_STACK_DEPTH {
//		if p.stack.Cur().hint.TokenEnd == "" || p.stack.Cur().hint.EndTokenOptional {
//			p.stack.Remove(p.buf)
//			p.buf.Reset()
//			//fmt.Println("removing from stack due to EOL/EOF")
//		} else {
//			break
//		}
//	}
//	// Return the buffer if the stack is empty
//	if p.stack.depth < MIN_STACK_DEPTH {
//		return line_element, nil
//	}
//	// Return syntax error if we didn't close all of our containers
//	return nil, fmt.Errorf("line %d: unexpected EOF while looking for token `%s'", p.stack.Cur().position.Line, p.stack.Cur().hint.TokenEnd)
//}

func (p *Parser) nextRune() (rune, error) {
	r, _, err := p.input.ReadRune()
	p.Offset++
	p.LineOffset++
	return r, err
}

func (p *Parser) unreadRune() error {
	err := p.input.UnreadRune()
	p.Offset--
	p.LineOffset--
	return err
}

func (p *Parser) nextLine() {
	p.Line++
	p.LineOffset = 0
}

//func (p *Parser) newElement() lang.Element {
//	e := lang.NewGeneric(p.stack.Cur().hint.Name)
//	//fmt.Printf("%#v\n", e)
//	if p.stack.Cur().hint.Factory != nil {
//		foo := p.stack.Cur().hint.Factory(e)
//		//fmt.Printf("%s\n", foo)
//		return foo
//	}
//	return e
//}
