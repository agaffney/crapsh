package lang

import (
	"fmt"
)

type Element interface {
	Output() string
}

type Generic struct {
	Line     uint
	Content  string
	children []Element
}

func NewGeneric(content string, line uint) *Generic {
	return &Generic{Content: content, Line: line}
}

func (g *Generic) Output() string {
	return g.Content
}

func (g *Generic) AddChild(e Element) {
	if g.children == nil {
		g.children = make([]Element, 0)
	}
	g.children = append(g.children, e)
}

func (g *Generic) String() string {
	return fmt.Sprintf("<Generic    content: %#v>", g.Content)
}

type FactoryFunc func(*Generic) Element

type ParserHint struct {
	Name            string
	TokenStart      string
	TokenEnd        string
	IgnoreEscapes   bool
	SkipEndToken    bool
	CaptureAll      bool
	EndOnEOF        bool
	EndOnWhitespace bool
	EndOnNewline    bool
	AllowedElements []string
	Factory         FactoryFunc
}

func GetElementHint(s string) *ParserHint {
	for _, foo := range ParserHints {
		if s == foo.Name {
			return foo
		}
	}
	return nil
}

func (h *ParserHint) AllowedElement(s string) bool {
	for _, foo := range h.AllowedElements {
		if s == foo {
			return true
		}
	}
	return false
}

var ParserHints []*ParserHint

func registerParserHints(hints []*ParserHint) {
	ParserHints = append(ParserHints, hints...)
}
