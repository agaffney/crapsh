package parser

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/agaffney/crapsh/parser/tree"
	"io"
	"strings"
	"unicode"
)

type Parser struct {
	input      *bufio.Reader
	Line       int
	Offset     int
	LineOffset int
}

func NewParser() *Parser {
	parser := &Parser{}
	return parser
}

func (p *Parser) Parse(input string) tree.Node {
	fmt.Printf("%#v\n", tree.Node_types)
	r := bufio.NewReader(strings.NewReader(input))
	p.input = r
	p.Line = 1
	p.LineOffset = 0
	p.Offset = 0
	topnode := tree.NewTop()
	p.Scan(topnode)
	return topnode
}

func (p *Parser) Scan(parent tree.Node) error {
	sn := tree.NewStatement(parent)
	var buf bytes.Buffer
	for {
		c, err := p.next_rune()
		if err != nil {
			if err != io.EOF {
				fmt.Printf("Error: %v\n", err)
			}
			break
		}
		fmt.Printf("Line %d, offset %d, overall offset %d: '%c' (%d)\n", p.Line, p.LineOffset, p.Offset, c, c)
		if unicode.IsSpace(c) || c == '\n' {
			fmt.Printf("buf contains: '%s'\n", buf.String())
			child := tree.NewGeneric(sn)
			child.Set_content(buf.String())
			sn.Add_child(child)
			buf.Reset()
			if c == '\n' {
				p.next_line()
			}
		} else {
			buf.WriteRune(c)
		}
	}
	fmt.Printf("%#v\n", parent)
	return nil
}

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
