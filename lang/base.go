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
			SkipCapture:      true,
			AllowedElements:  []string{"Command"},
		},
		{
			Name:             `Command`,
			TokenEnd:         tokens.SEMICOLON,
			EndTokenOptional: true,
			SkipCapture:      true,
			AllowedElements:  []string{"Argument"},
		},
		{
			Name:            `Argument`,
			EndOnWhitespace: true,
			//SkipCapture:     true,
			AllowedElements: []string{"StringSingle", "StringDouble", "DollarSign", "Generic"},
		},
		{
			Name:       `Generic`,
			CaptureAll: true,
		},
		{
			Name:            `DollarSign`,
			TokenStart:      tokens.DOLLAR_SIGN,
			SkipCapture:     true,
			AllowedElements: []string{"Variable"},
		},
	})

}
