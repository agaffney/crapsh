package lang

import (
	"github.com/agaffney/crapsh/parser/tokens"
)

func init() {
	registerParserHints([]*ParserHint{
		{
			Name:            `Line`,
			TokenEnd:        tokens.NEWLINE,
			EndOnEOF:        true,
			AllowedElements: []string{"Command"},
		},
		{
			Name:            `Command`,
			TokenEnd:        tokens.SEMICOLON,
			EndOnNewline:    true,
			EndOnEOF:        true,
			AllowedElements: []string{"Argument"},
		},
		{
			Name:            `Argument`,
			EndOnWhitespace: true,
			EndOnNewline:    true,
			EndOnEOF:        true,
			AllowedElements: []string{"StringSingle", "StringDouble", "Variable"},
		},
	})

}
