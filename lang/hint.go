package lang

const (
	HINT_TYPE_NODE = iota
	HINT_TYPE_ELEMENT
	HINT_TYPE_TOKEN
	HINT_TYPE_GROUP
	HINT_TYPE_ANY
)

type ParserHint struct {
	Name     string
	Type     int
	Optional bool
	Many     bool
	Members  []*ParserHint
}
