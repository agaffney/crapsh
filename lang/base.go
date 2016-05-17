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

type Base struct {
	Line    uint
	Content string
}

type Generic struct {
	*Base
}

type FactoryFunc func(*Base) Element

func New(content string, line uint) *Base {
	return &Base{Content: content, Line: line}
}

func NewGeneric(base *Base) Element {
	return &Generic{Base: base}
}

func (b *Base) Output() string {
	return b.Content
}

func (c *ContainerBase) AddChild(e Element) {
	if c.children == nil {
		c.children = make([]Element, 0)
	}
	c.children = append(c.children, e)
}
