package lang

const (
	STRING_TYPE_INVALID int = iota
	STRING_TYPE_SINGLE
	STRING_TYPE_DOUBLE
)

type String struct {
	*Base
	*ContainerBase
	Type int
}

func NewStringSingle(base *Base) Element {
	return &String{Base: base, Type: STRING_TYPE_SINGLE}
}

func NewStringDouble(base *Base) Element {
	return &String{Base: base, Type: STRING_TYPE_DOUBLE}
}
