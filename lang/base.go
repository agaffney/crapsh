package lang

type Element interface {
	Output() string
}

type Container interface {
	AddChild(Element)
}

type Generic struct {
	Line     uint
	Content  string
	children []Element
}

type FactoryFunc func(*Generic) Element

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
