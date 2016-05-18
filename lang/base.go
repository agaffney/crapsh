package lang

type Element interface {
	Output() string
}

type Container interface {
	AddChild(Element)
}

type ContainerBase struct {
	children []Element
}

type Generic struct {
	Line    uint
	Content string
}

type FactoryFunc func(*Generic) Element

func NewGeneric(content string, line uint) *Generic {
	return &Generic{Content: content, Line: line}
}

func (g *Generic) Output() string {
	return g.Content
}

func (c *ContainerBase) AddChild(e Element) {
	if c.children == nil {
		c.children = make([]Element, 0)
	}
	c.children = append(c.children, e)
}
