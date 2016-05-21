package lang

import (
	"github.com/agaffney/crapsh/parser/tokens"
)

func init() {
	registerParserHints([]*ParserHint{
		{
			Name:             `Line`,
			TokenEnd:         tokens.NEWLINE,
			EndTokenOptional: true,
			AllowedElements:  []string{"Command"},
		},
		{
			Name:             `Command`,
			TokenEnd:         tokens.SEMICOLON,
			EndTokenOptional: true,
			AllowedElements:  []string{"Argument"},
		},
		{
			Name:            `Argument`,
			EndOnWhitespace: true,
			AllowedElements: []string{"StringSingle", "StringDouble", "Variable"},
		},
	})

}
