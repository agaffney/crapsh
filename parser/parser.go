package parser

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/agaffney/crapsh/lang"
	"io"
	"strings"
	"unicode"
)

const (
	MIN_STACK_DEPTH = 0
)

type Parser struct {
	Position
	input      *bufio.Reader
	stack      []*StackEntry
	stackdepth int
}

type Position struct {
	Line       uint
	Offset     uint
	LineOffset uint
}

type StackEntry struct {
	hint     *lang.ParserHint
	allowed  []*lang.ParserHint
	position Position
	element  lang.Element
}

func NewParser() *Parser {
	parser := &Parser{}
	return parser
}

func (p *Parser) Parse(input string) {
	r := bufio.NewReader(strings.NewReader(input))
	p.input = r
	p.Line = 1
	p.LineOffset = 0
	p.Offset = 0
	for {
		line, err := p.GetNextLine()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			break
		}
		// EOF
		if line.NumChildren() == 0 {
			break
		}
		fmt.Printf("Line: %s\n", line)
	}
}

func (p *Parser) GetNextLine() (lang.Element, error) {
	escape := false
	buf := bytes.NewBuffer(nil)
	// Reset the hint stack
	p.stack = nil
	p.stackAdd(lang.GetElementHints([]string{"Line"})[0])
	line_element := p.stackCur().element
	for {
		c, err := p.nextRune()
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			break
		}
		//fmt.Printf("Stack item (%d): %#v\n", p.stackdepth+1, p.stack[p.stackdepth])
		fmt.Printf("Line %d, offset %d, overall offset %d: %#U, buf: %#v\n", p.Line, p.LineOffset, p.Offset, c, buf.String())
		if c == '\\' && !p.stack[p.stackdepth].hint.IgnoreEscapes {
			escape = !escape
			// Explicitly skip to the next iteration so we don't hit
			// the code below to turn off the 'escape' flag
			continue
		} else {
			if c == '\n' {
				p.nextLine()
				if escape {
					buf.WriteRune(c)
				} else {
					if p.stack[p.stackdepth].hint.TokenEnd == "" || p.stack[p.stackdepth].hint.EndTokenOptional {
						p.stackRemove(buf)
						fmt.Println("removing from stack due to newline")
						buf.Reset()
						break
					}
				}
				continue
			}
			if unicode.IsSpace(c) && !escape && p.stackCur().hint.EndOnWhitespace {
				if p.stackCur().hint.CaptureAll {
					// We're using the EndOnWhitespace value from our "parent", so if it's found,
					// we should remove the CaptureAll element from the stack
					p.stackRemove(buf)
					fmt.Println("removing from stack due to CaptureAll and whitespace")
					buf.Reset()
				}
				p.stackRemove(buf)
				fmt.Println("removing from stack due to whitespace")
				buf.Reset()
				continue
			}
			// Add new character to tmpbuf for checking start/end tokens
			//tmpbuf.Reset()
			//tmpbuf.Write(buf.Bytes())
			//tmpbuf.WriteRune(c)
			buf.WriteRune(c)
			if escape == false && p.stackCur().hint.TokenEnd != "" && checkBufForToken(buf, p.stack[p.stackdepth].hint.TokenEnd) {
				// Remove start token from buf
				buf = bytes.NewBuffer(buf.Bytes()[:buf.Len()-len(p.stackCur().hint.TokenEnd)])
				if p.stackCur().hint.CaptureAll {
					// We're using the end token from our "parent", so if it's found,
					// we should remove the CaptureAll element from the stack
					p.stackRemove(buf)
					fmt.Println("removing from stack due to CaptureAll and finding end token")
					buf.Reset()
				}
				p.stackRemove(buf)
				fmt.Println("removing from stack due to finding end token")
				buf.Reset()
				continue
			}
			found := false
			for _, hint := range p.stackCur().allowed {
				if checkBufForToken(buf, hint.TokenStart) { //|| hint.CaptureAll {
					// Remove start token from buf
					buf = bytes.NewBuffer(buf.Bytes()[:buf.Len()-len(hint.TokenStart)])
					if p.stackCur().hint.CaptureAll {
						// We're using the allowed elements from our "parent", so if one is found,
						// we should remove the CaptureAll element from the stack
						p.stackRemove(buf)
						fmt.Println("removing from stack due to CaptureAll and finding start token")
						buf.Reset()
					}
					p.stackAdd(hint)
					if hint.SkipCapture {
						p.unreadRune()
						buf.Reset()
					}
					found = true
					break
				}
			}
			if !found && !p.stackCur().hint.CaptureAll {
				return nil, fmt.Errorf("line %d, pos %d: unexpected character `%c'", p.Position.Line, p.Position.LineOffset, c)
			}
		}
		// Reset the 'escape' flag
		escape = false
	}
	// Remove any stack items that allow ending on EOF
	for p.stackdepth >= MIN_STACK_DEPTH {
		if p.stack[p.stackdepth].hint.TokenEnd == "" || p.stack[p.stackdepth].hint.EndTokenOptional {
			p.stackRemove(buf)
			buf.Reset()
			fmt.Println("removing from stack due to EOL/EOF")
		} else {
			break
		}
	}
	// Return the buffer if the stack is empty
	if p.stackdepth < MIN_STACK_DEPTH {
		return line_element, nil
	}
	// Return syntax error if we didn't close all of our containers
	return nil, fmt.Errorf("line %d: unexpected EOF while looking for token `%s'", p.stack[p.stackdepth].position.Line, p.stack[p.stackdepth].hint.TokenEnd)
}

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

func (p *Parser) newElement() lang.Element {
	e := lang.NewGeneric("", p.Line, p.stackCur().hint.Name)
	//fmt.Printf("%#v\n", e)
	if p.stack[p.stackdepth].hint.Factory != nil {
		foo := p.stack[p.stackdepth].hint.Factory(e)
		//fmt.Printf("%s\n", foo)
		return foo
	}
	return e
}

func (p *Parser) stackAdd(hint *lang.ParserHint) {
	if p.stack == nil {
		p.stack = make([]*StackEntry, 0)
		p.stackdepth = -1
	}
	allowed := []*lang.ParserHint{}
	if hint.CaptureAll {
		// If the current stack entry captures all, we need to use the allowed
		// elements from the "parent". We also want to filter ourselves out
		for _, foo := range lang.GetElementHints(p.stackCur().hint.AllowedElements) {
			if !foo.CaptureAll {
				allowed = append(allowed, foo)
			}
		}
		// Use the end token info from the parent
		hint.TokenEnd = p.stackCur().hint.TokenEnd
		hint.EndOnWhitespace = p.stackCur().hint.EndOnWhitespace
	} else {
		allowed = lang.GetElementHints(hint.AllowedElements)
	}
	e := &StackEntry{hint, allowed, p.Position, nil}
	p.stack = append(p.stack, e)
	p.stackdepth++
	e.element = p.newElement()
	fmt.Printf("\nstack[%d] = %#v\n\n", p.stackdepth, hint)
	//fmt.Printf("  allowed = [\n")
	//for _, foo := range e.allowed {
	//	fmt.Printf("    %#v,\n", foo)
	//}
	//fmt.Printf("  ]\n\n")
}

func (p *Parser) stackRemove(buf *bytes.Buffer) {
	p.stackCur().element.SetContent(buf.String())
	if p.stackdepth > MIN_STACK_DEPTH {
		p.stackPrev().element.AddChild(p.stackCur().element)
	}
	p.stack = p.stack[:len(p.stack)-1]
	p.stackdepth--
	if p.stackdepth >= MIN_STACK_DEPTH {
		fmt.Printf("\nstack[%d] = %#v\n\n", p.stackdepth, p.stack[p.stackdepth].hint)
	}
}

func (p *Parser) stackCur() *StackEntry {
	if p.stackdepth >= MIN_STACK_DEPTH {
		return p.stack[p.stackdepth]
	}
	return nil
}

func (p *Parser) stackPrev() *StackEntry {
	if p.stackdepth-1 >= MIN_STACK_DEPTH {
		return p.stack[p.stackdepth-1]
	}
	return nil
}

// Grab n bytes (length of token) from end of buf and compare to token
func checkBufForToken(buf *bytes.Buffer, token string) bool {
	token_len := len(token)
	if buf.Len() < token_len {
		return false
	}
	buf_bytes := buf.Bytes()[buf.Len()-token_len:]
	for i, b := range []byte(token) {
		if buf_bytes[i] != b {
			return false
		}
	}
	return true
}
