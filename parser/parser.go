package parser

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
	"unicode"
)

type Parser struct {
	Position
	input *bufio.Reader
}

type Position struct {
	Line       int
	Offset     int
	LineOffset int
}

type ContainerStackEntry struct {
	container *Container
	position  Position
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
		line, err := p.Get_next_line()
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

func (p *Parser) Get_next_line() (*bytes.Buffer, error) {
	var buf bytes.Buffer
	var escape = false
	var stackdepth int = 0
	stack := []ContainerStackEntry{ContainerStackEntry{line_container, p.Position}}
	for {
		c, err := p.next_rune()
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			break
		}
		//fmt.Printf("Stack item (%d): %#v\n", stackdepth+1, stack[stackdepth])
		fmt.Printf("Line %d, offset %d, overall offset %d: %#U\n", p.Line, p.LineOffset, p.Offset, c)
		if c == '\\' {
			escape = true
			// Explicitly skip to the next iteration so we don't hit
			// the code below to turn off the 'escape' flag
			continue
		} else {
			buf.WriteRune(c)
			if escape == false && check_buf_for_token(&buf, stack[stackdepth].container.TokenEnd) {
				stack = stack[:len(stack)-1]
				stackdepth--
			} else if stack[stackdepth].container.AllowedContainers != nil {
				for _, cont := range containers {
					if stack[stackdepth].container.Allowed_container(cont.Name) {
						//fmt.Printf("%#v\n", cont)
						if check_buf_for_token(&buf, cont.Token) {
							stack = append(stack, ContainerStackEntry{cont, p.Position})
							stackdepth++
							break
						}
					}
				}
			}
			if c == '\n' {
				p.next_line()
			}
		}
		if stackdepth < 0 {
			return &buf, nil
		}
		// Reset the 'escape' flag
		escape = false
	}
	return nil, fmt.Errorf("line %d: unexpected EOF while looking for token `%s'", stack[stackdepth].position.Line, stack[stackdepth].container.TokenEnd)
}

func (p *Parser) Scan() error {
	var buf bytes.Buffer
	for {
		c, err := p.next_rune()
		if err != nil {
			if err != io.EOF {
				fmt.Printf("Error: %v\n", err)
			}
			//			p.process_scan_buf(&buf, sn)
			break
		}
		fmt.Printf("Line %d, offset %d, overall offset %d: %#U\n", p.Line, p.LineOffset, p.Offset, c)
		if unicode.IsSpace(c) || c == '\n' {
			//			p.process_scan_buf(&buf, sn)
		} else if c == '\n' {
			//			p.process_scan_buf(&buf, sn)
			p.next_line()
		} else {
			buf.WriteRune(c)
		}
	}
	return nil
}

//func (p *Parser) process_scan_buf(buf *bytes.Buffer, parent tree.Node) {
//	fmt.Printf("buf contains: '%s'\n", buf.String())
//	child := tree.NewGeneric(parent)
//	child.Set_content(buf.String())
//	parent.Add_child(child)
//	buf.Reset()
//}

func (p *Parser) next_rune() (rune, error) {
	r, _, err := p.input.ReadRune()
	p.Offset++
	p.LineOffset++
	return r, err
}

func (p *Parser) unread_rune() error {
	err := p.input.UnreadRune()
	p.Offset--
	p.LineOffset--
	return err
}

func (p *Parser) next_line() {
	p.Line++
	p.LineOffset = 0
}

// Grab n bytes (length of token) from end of buf and compare to token
func check_buf_for_token(buf *bytes.Buffer, token string) bool {
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
