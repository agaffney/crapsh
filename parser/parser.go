package parser

import (
	//"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/agaffney/crapsh/lang"
	parser_input "github.com/agaffney/crapsh/parser/input"
	"github.com/agaffney/crapsh/parser/lexer"
	//"github.com/agaffney/crapsh/util"
	"io"
	"log"
)

var ERR_FOUND_PARENT_END_TOKEN = errors.New("found parent end token")

type Parser struct {
	//input       *bufio.Reader
	input       parser_input.Input
	stack       *Stack
	buf         *bytes.Buffer
	tokenBuf    []*lexer.Token
	tokenIdx    int
	commandChan chan lang.Element
	errorChan   chan error
	lexer       *lexer.Lexer
}

func NewParser() *Parser {
	parser := &Parser{}
	parser.lexer = lexer.New()
	return parser
}

func (p *Parser) Reset() {
	p.commandChan = make(chan lang.Element)
	p.errorChan = make(chan error)
	p.stack = &Stack{parser: p}
	p.buf = bytes.NewBuffer(nil)
	p.lexer.Reset()
}

func (p *Parser) Parse(input parser_input.Input) {
	p.Reset()
	p.lexer.Start(input)
	for {
		token, err := p.lexer.NextToken()
		if err != nil {
			if err != io.EOF {
				log.Fatal(err)
			}
		}
		if token == nil {
			break
		}
		fmt.Printf("token = %#v\n", token)
	}
	return
	/*
		p.tokenBuf = make([]*lexer.Token, 0)
		p.tokenIdx = -1
		go func() {
			for {
				// Reset the hint stack
				p.stack.Reset()
				// Start parsing with "root" element
				ok, err := p.parseElement("Root", true)
				if err != nil {
					p.errorChan <- err
					close(p.errorChan)
					break
				}
				// EOF
				if !ok {
					break
				}
			}
			close(p.commandChan)
		}()
	*/
}

/*
func (p *Parser) GetError() error {
	select {
	case err := <-p.errorChan:
		return err
	default:
		return nil
	}
}
*/

func (p *Parser) GetCommand() *lang.Element {
	cmd := <-p.commandChan
	if cmd != nil {
		return &cmd
	} else {
		return nil
	}
}

/*
// Handles an individual parser hint
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
			ok, err = p.parseElement(hint.Name, false)
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
			fmt.Printf("parseHandleHint(): !ok, hint.Many=%t, count=%d\n", hint.Many, count)
			p.setTokenIdx(origTokenIdx)
			if hint.Optional || (hint.Many && count > 0) {
				return true, nil
			}
			return false, nil
		}
		if !hint.Many {
			break
		}
		count++
		// TODO: remove look ahead, as we're being specific about valid tokens in the hints
		// Look ahead to see if next token in buffer matches current end token
		if nextHint := p.stack.Cur().NextHint(); nextHint != nil && nextHint.Type == lang.HINT_TYPE_TOKEN {
			token, err := p.nextToken()
			if err != nil {
				if err == io.EOF {
					return ok, nil
				}
				return ok, err
			}
			if token != nil && nextHint.Name == token.Type {
				fmt.Printf("parseHandleHint(): look-ahead MATCH: token=%#v, nextHint=%#v\n", token, nextHint)
				p.prevToken()
				if hint.Optional || (hint.Many && count > 0) {
					return true, nil
				}
				return ok, nil
			}
			p.prevToken()
		}
	}
	return true, nil
}

func (p *Parser) parseToken(hint *lang.ParserHint) (bool, error) {
	token, err := p.nextToken()
	if err != nil {
		if err == io.EOF {
			return false, nil
		}
		return false, err
	}
	util.DumpObject(hint, "parseToken(): hint = ")
	util.DumpObject(p.curToken(), "parseToken(): curToken = ")
	// TODO: remove parentEndToken logic, as it's not needed
	if nextHint := p.stack.Cur().NextHint(); nextHint != nil {
		if nextHint.Type == lang.HINT_TYPE_TOKEN && nextHint.Name == token.Type {
			p.prevToken()
			return false, nil
		}
	} else {
		if p.stack.Cur().parentEndToken != nil && p.stack.Cur().parentEndToken.Name == token.Type {
			fmt.Printf("parseToken(): matched parentTokenEnd=%#v\n", p.stack.Cur().parentEndToken)
			p.prevToken()
			return false, ERR_FOUND_PARENT_END_TOKEN
		}
	}
	tokenMatch := false
	for _, hint_token := range hint.Tokens {
		if hint_token == token.Type {
			tokenMatch = true
			break
		}
	}
	if tokenMatch {
		e := lang.NewGeneric(`Token`)
		e.Content = token.Value
		p.stack.Cur().element.AddChild(e)
		return true, nil
	}
	return false, nil
}

// Handles a 'group' parser hint
// Succeeds if all parser hints match
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

// Handles an 'element' parser hint
func (p *Parser) parseElement(element string, send_channel bool) (bool, error) {
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
		if parentEndToken == nil && p.stack.Cur().parentEndToken != nil {
			parentEndToken = p.stack.Cur().parentEndToken
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
		if err != ERR_FOUND_PARENT_END_TOKEN {
			return false, err
		} else {
			// If the parent end token has been found and matches
			// our current parent end token, pass the error up the
			// chain
			if p.stack.Cur().parentEndToken.Name == p.curToken().Type {
				return false, err
			}
		}
	}
	p.stack.Remove()
	if ok {
		if p.stack.Cur() != nil {
			p.stack.Cur().element.AddChild(e)
		} else if send_channel {
			//util.DumpJson(e, "sending root element:\n")
			p.commandChan <- e
		}
	}
	return ok, nil
}

// Handles an 'any' parser hint
// Succeeds if any parser hints match
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
*/

func (p *Parser) getTokenIdx() int {
	return p.tokenIdx
}

func (p *Parser) setTokenIdx(idx int) {
	p.tokenIdx = idx
	if idx >= 0 {
		fmt.Printf("setTokenIdx(%d): curToken = %#v\n", idx, p.curToken())
	}
}

func (p *Parser) curToken() *lexer.Token {
	return p.tokenBuf[p.tokenIdx]
}

func (p *Parser) prevToken() {
	if p.tokenIdx > 0 {
		p.tokenIdx--
	}
}

func (p *Parser) nextToken() (*lexer.Token, error) {
	if p.tokenIdx < len(p.tokenBuf)-1 {
		p.tokenIdx++
		return p.curToken(), nil
	} else {
		token := p.lexer.ReadToken()
		/*
			if err != nil {
				return nil, err
			}
		*/
		if token == nil {
			return nil, io.EOF
		}
		p.tokenBuf = append(p.tokenBuf, token)
		p.tokenIdx++
		return token, nil
	}
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
