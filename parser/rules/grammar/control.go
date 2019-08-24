package grammar

import (
	"github.com/agaffney/crapsh/parser/tokens"
)

func init() {
	registerRules([]*GrammarRule{
		{
			Name: `if_clause`,
			ParserHints: []*ParserHint{
				{
					Type:       HINT_TYPE_TOKEN,
					TokenTypes: []int{tokens.TOKEN_IF},
				},
			},
		},
	})
}
