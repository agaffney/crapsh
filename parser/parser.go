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
			line, err := p.parseElement("Line")
			util.DumpJson(line, "Line: ")
			if err != nil {
				p.Error = err
				break
			}
			// EOF
			if line == nil || line.NumChildren() == 0 {
				break
			}
			p.LineChan <- line
		}
		close(p.LineChan)
	}()
}

func (p *Parser) parseHandleHint(hint *lang.ParserHint) (lang.Element, error) {
	//util.DumpObject(hint, "parseHandleHint() hint=")
	var foo lang.Element
	var err error
	var count int
	for {
		switch {
		case hint.Type == lang.HINT_TYPE_ELEMENT:
			foo, err = p.parseElement(hint.Name)
		case hint.Type == lang.HINT_TYPE_ANY:
			foo, err = p.parseAny(hint.Members)
		case hint.Type == lang.HINT_TYPE_GROUP:
			foo, err = p.parseGroup(hint.Members)
		case hint.Type == lang.HINT_TYPE_TOKEN:
			foo, err = p.parseToken(hint)
		default:
			return nil, fmt.Errorf("Unhandled hint type: %d\n", hint.Type)
		}
		if err != nil {
			return nil, err
		}
		if foo == nil {
			if hint.Optional {
				return nil, nil
			}
			if hint.Many && count > 0 {
				return nil, nil
			}
			return nil, nil
		}
		p.stack.Prev().element.AddChild(foo)
		if !hint.Many {
			break
		}
		count++
	}
	return nil, nil
}

func (p *Parser) parseToken(hint *lang.ParserHint) (lang.Element, error) {
	token, err := p.nextToken()
	if err != nil {
		return nil, err
	}
	if token == nil {
		return nil, nil
	}
	util.DumpObject(token, "parseToken() token=")
	p.tokenBuf = append(p.tokenBuf, token)
	p.tokenIdx++
	foo := p.curToken()
	if hint.Name == foo.Type {
		e := lang.NewGeneric(p.stack.Cur().entry.Name)
		e.Content = foo.Value
		return e, nil
	}
	return nil, nil
}

func (p *Parser) parseGroup(hints []*lang.ParserHint) (lang.Element, error) {
	for _, hint := range hints {
		util.DumpObject(hint, "parseGroup(): hint = ")
		_, err := p.parseHandleHint(hint)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (p *Parser) parseElement(element string) (lang.Element, error) {
	util.DumpObject(element, "parseElement() element=")
	entry := lang.GetElementEntry(element)
	if entry == nil {
		return nil, nil
	}
	p.stack.Add(entry)
	e := lang.NewGeneric(p.stack.Cur().entry.Name)
	//fmt.Printf("%#v\n", e)
	if p.stack.Cur().entry.Factory != nil {
		foo := p.stack.Cur().entry.Factory(e)
		//fmt.Printf("%s\n", foo)
		p.stack.Cur().element = foo
	}
	p.stack.Cur().element = e
	_, err := p.parseGroup(entry.ParserData)
	if err != nil {
		return nil, err
	}
	p.stack.Remove()
	return nil, nil
}

func (p *Parser) parseAny(hints []*lang.ParserHint) (lang.Element, error) {
	for _, hint := range hints {
		util.DumpObject(hint, "parseAny(): hint = ")
		foo, err := p.parseHandleHint(hint)
		if err != nil {
			return nil, err
		}
		if foo != nil {
			return foo, nil
		}
	}
	return nil, nil
}

func (p *Parser) scanTokens() (lang.Element, error) {
	// Reset the hint stack
	p.stack.Reset()
	p.stack.Add(lang.GetElementEntry("Line"))
	line_element := p.stack.Cur().element
	for {
		token, err := p.nextToken()
		if err != nil {
			return nil, err
		}
		if token == nil {
			return nil, nil
		}
		fmt.Printf("Token: %#v\n", token)
		p.tokenBuf = append(p.tokenBuf, token)
		util.DumpJson(p.tokenBuf, "scanTokens() tokenBuf=")
	}
	return line_element, nil
}

func (p *Parser) getTokenIdx() int {
	return p.tokenIdx
}

func (p *Parser) setTokenIdx(idx int) {
	p.tokenIdx = idx
}

func (p *Parser) curToken() *Token {
	return p.tokenBuf[p.tokenIdx]
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
