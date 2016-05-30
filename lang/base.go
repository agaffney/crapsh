package lang

import (
	"github.com/agaffney/crapsh/lang/tokens"
)

func init() {
	registerParserHints([]*ParserHint{
		{
			Name:            `Line`,
			SkipCapture:     true,
			AllowedElements: []string{"Command"},
		},
		{
			Name:             `Command`,
			TokenEnd:         tokens.SEMICOLON,
			EndTokenOptional: true,
			SkipCapture:      true,
			AllowedElements:  []string{"Subshell", "Argument"},
		},
		{
			Name:            `Argument`,
			EndOnWhitespace: true,
			SkipCapture:     true,
			AllowedElements: []string{"StringSingle", "StringDouble", "SubshellCapture", "SubshellBacktick", "Generic"},
		},
		{
			Name:       `Generic`,
			CaptureAll: true,
		},
	})

}
