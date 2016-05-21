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
	Name             string
	TokenStart       string
	TokenEnd         string
	IgnoreEscapes    bool
	EndOnWhitespace  bool
	EndTokenOptional bool
	SkipCapture      bool
	CaptureAll       bool
	AllowedElements  []string
	Factory          FactoryFunc
}

func GetElementHints(elements []string) []*ParserHint {
	hints := []*ParserHint{}
	for _, e := range elements {
		for _, hint := range ParserHints {
			if e == hint.Name {
				hints = append(hints, hint)
				break
			}
		}
	}
	return hints
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
