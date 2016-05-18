package lang

const (
	STRING_TYPE_INVALID int = iota
	STRING_TYPE_SINGLE
	STRING_TYPE_DOUBLE
)

type String struct {
	*Generic
	Type int
}

func NewStringSingle(base *Generic) Element {
	return &String{Generic: base, Type: STRING_TYPE_SINGLE}
}

func NewStringDouble(base *Generic) Element {
	return &String{Generic: base, Type: STRING_TYPE_DOUBLE}
}
