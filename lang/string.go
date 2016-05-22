package lang

import (
	"github.com/agaffney/crapsh/parser/tokens"
)

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

func init() {
	registerParserHints([]*ParserHint{
		{
			Name:            `StringSingle`,
			TokenStart:      tokens.SINGLE_QUOTE,
			TokenEnd:        tokens.SINGLE_QUOTE,
			IgnoreEscapes:   true,
			AllowedElements: []string{"Generic"},
			Factory:         NewStringSingle,
		},
		{
			Name:            `StringDouble`,
			TokenStart:      tokens.DOUBLE_QUOTE,
			TokenEnd:        tokens.DOUBLE_QUOTE,
			AllowedElements: []string{"Variable", "SubshellCapture", "Subshell", "SubshellBacktick", "Generic"},
			Factory:         NewStringDouble,
		},
	})
}
