package parser

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/agaffney/crapsh/lang"
	"io"
	"strings"
)

type Parser struct {
	Position
	input      *bufio.Reader
	stack      []*HintStackEntry
	stackdepth int
}

type Position struct {
	Line       uint
	Offset     uint
	LineOffset uint
}

type HintStackEntry struct {
	hint     *lang.ParserHint
	allowed  []*lang.ParserHint
	position Position
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
		if line.Len() == 0 {
			break
		}
		fmt.Printf("Line: %s\n", line)
	}
}

func (p *Parser) GetNextLine() (*bytes.Buffer, error) {
	var buf bytes.Buffer
	var linebuf bytes.Buffer
	var escape = false
	// Reset the hint stack
	p.stack = nil
	p.stackAdd(lang.GetElementHints([]string{"Line"})[0])
	for {
		c, err := p.nextRune()
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			break
		}
		//fmt.Printf("Stack item (%d): %#v\n", p.stackdepth+1, p.stack[p.stackdepth])
		fmt.Printf("Line %d, offset %d, overall offset %d: %#U\n", p.Line, p.LineOffset, p.Offset, c)
		if c == '\\' && !p.stack[p.stackdepth].hint.IgnoreEscapes {
			escape = !escape
			// Explicitly skip to the next iteration so we don't hit
			// the code below to turn off the 'escape' flag
			continue
		} else {
			buf.WriteRune(c)
			if escape == false && checkBufForToken(&buf, p.stack[p.stackdepth].hint.TokenEnd) {
				if p.stack[p.stackdepth].hint.Factory != nil {
					foo := p.stack[p.stackdepth].hint.Factory(lang.NewGeneric(buf.String(), p.Line))
					fmt.Printf("%s\n", foo)
				}
				p.stackRemove()
				linebuf.Write(buf.Bytes())
				buf.Reset()
			} else {
				for _, hint := range p.stackLast().allowed {
					if checkBufForToken(&buf, hint.TokenStart) {
						p.stackAdd(hint)
						break
					}
				}
			}
			if c == '\n' {
				p.nextLine()
				if p.stack[p.stackdepth].hint.TokenEnd == "" || p.stack[p.stackdepth].hint.EndTokenOptional {
					p.stackRemove()
					linebuf.Write(buf.Bytes())
					buf.Reset()
					break
				}
			}
		}
		// Reset the 'escape' flag
		escape = false
	}
	// Remove any stack items that allow ending on EOF
	for p.stackdepth >= 0 {
		if p.stack[p.stackdepth].hint.TokenEnd == "" || p.stack[p.stackdepth].hint.EndTokenOptional {
			p.stackRemove()
		}
	}
	// Return the buffer if the stack is empty
	if p.stackdepth < 0 {
		return &linebuf, nil
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

func (p *Parser) stackAdd(hint *lang.ParserHint) {
	if p.stack == nil {
		p.stack = make([]*HintStackEntry, 0)
		p.stackdepth = -1
	}
	e := &HintStackEntry{hint, lang.GetElementHints(hint.AllowedElements), p.Position}
	p.stack = append(p.stack, e)
	p.stackdepth++
	//fmt.Printf("\nstack[%d] = %#v\n", p.stackdepth, hint)
	//fmt.Printf("  allowed = [\n")
	//for _, foo := range e.allowed {
	//	fmt.Printf("    %#v,\n", foo)
	//}
	//fmt.Printf("  ]\n\n")
}

func (p *Parser) stackRemove() {
	p.stack = p.stack[:len(p.stack)-1]
	p.stackdepth--
}

func (p *Parser) stackLast() *HintStackEntry {
	if p.stackdepth >= 0 {
		return p.stack[p.stackdepth]
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
