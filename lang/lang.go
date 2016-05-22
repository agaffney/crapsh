package lang

import (
	"fmt"
)

type Element interface {
	Output() string
	AddChild(Element)
	SetContent(string)
	NumChildren() int
}

type Generic struct {
	Content  string
	Children []Element
	TypeName string
}

func NewGeneric(typeName string) *Generic {
	return &Generic{TypeName: typeName}
}

func (g *Generic) Output() string {
	return g.Content
}

func (g *Generic) AddChild(e Element) {
	if g.Children == nil {
		g.Children = make([]Element, 0)
	}
	g.Children = append(g.Children, e)
	//fmt.Printf("new child: %s\n", e)
}

func (g *Generic) SetContent(content string) {
	g.Content = content
}

func (g *Generic) NumChildren() int {
	return len(g.Children)
}

func (g *Generic) String() string {
	typeName := g.TypeName
	if typeName == "" {
		typeName = "Generic"
	}
	return fmt.Sprintf("<%s  content: %#v, children: %#v>", typeName, g.Content, g.Children)
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
