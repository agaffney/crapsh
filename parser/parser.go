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
	input *bufio.Reader
}

type Position struct {
	Line       uint
	Offset     uint
	LineOffset uint
}

type HintStackEntry struct {
	hint     *lang.ParserHint
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
	var stackdepth int = 0
	stack := []HintStackEntry{HintStackEntry{lang.GetElementHint("Line"), p.Position}}
	for {
		c, err := p.nextRune()
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			break
		}
		//fmt.Printf("Stack item (%d): %#v\n", stackdepth+1, stack[stackdepth])
		fmt.Printf("Line %d, offset %d, overall offset %d: %#U\n", p.Line, p.LineOffset, p.Offset, c)
		if c == '\\' && !stack[stackdepth].hint.IgnoreEscapes {
			escape = true
			// Explicitly skip to the next iteration so we don't hit
			// the code below to turn off the 'escape' flag
			continue
		} else {
			buf.WriteRune(c)
			if escape == false && checkBufForToken(&buf, stack[stackdepth].hint.TokenEnd) {
				if stack[stackdepth].hint.Factory != nil {
					foo := stack[stackdepth].hint.Factory(lang.NewGeneric(buf.String(), p.Line))
					fmt.Printf("%s\n", foo)
				}
				stack = stack[:len(stack)-1]
				stackdepth--
				linebuf.Write(buf.Bytes())
				buf.Reset()
			} else if stack[stackdepth].hint.AllowedElements != nil {
				for _, cont := range lang.ParserHints {
					if stack[stackdepth].hint.AllowedElement(cont.Name) {
						//fmt.Printf("%#v\n", cont)
						if checkBufForToken(&buf, cont.TokenStart) {
							stack = append(stack, HintStackEntry{cont, p.Position})
							stackdepth++
							break
						}
					}
				}
			}
			if c == '\n' {
				p.nextLine()
				if stack[stackdepth].hint.EndOnNewline {
					stack = stack[:len(stack)-1]
					stackdepth--
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
	for stackdepth >= 0 {
		if stack[stackdepth].hint.EndOnEOF {
			stack = stack[:len(stack)-1]
			stackdepth--
		}
	}
	// Return the buffer if the stack is empty
	if stackdepth < 0 {
		return &linebuf, nil
	}
	// Return syntax error if we didn't close all of our containers
	return nil, fmt.Errorf("line %d: unexpected EOF while looking for token `%s'", stack[stackdepth].position.Line, stack[stackdepth].hint.TokenEnd)
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
