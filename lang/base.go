package lang

import (
	"github.com/agaffney/crapsh/parser/tokens"
)

func init() {
	registerParserHints([]*ParserHint{
		{
			Name:              `Line`,
			TokenEnd:          tokens.NEWLINE,
			AllowEndOnEOF:     true,
			AllowedContainers: []string{"StringSingle", "StringDouble", "Variable"},
		},
	})

}
