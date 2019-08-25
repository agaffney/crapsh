package grammar

import (
	"github.com/agaffney/crapsh/parser/tokens"
)

func init() {
	registerRules([]*GrammarRule{
		{
			Name: `io_redirect`,
			ParserHints: []*ParserHint{
				{
					Type:       HINT_TYPE_TOKEN,
					Optional:   true,
					TokenTypes: []int{tokens.TOKEN_IO_NUMBER},
				},
				{
					Type: HINT_TYPE_ANY,
					Members: []*ParserHint{
						{
							Type:     HINT_TYPE_RULE,
							RuleName: `io_file`,
						},
						{
							Type:     HINT_TYPE_RULE,
							RuleName: `io_here`,
						},
					},
				},
			},
		},
		{
			Name: `io_file`,
			ParserHints: []*ParserHint{
				{
					Type: HINT_TYPE_TOKEN,
					TokenTypes: []int{
						tokens.TOKEN_LESS,
						tokens.TOKEN_LESSAND,
						tokens.TOKEN_GREAT,
						tokens.TOKEN_GREATAND,
						tokens.TOKEN_DGREAT,
						tokens.TOKEN_LESSGREAT,
						tokens.TOKEN_CLOBBER,
					},
				},
				{
					Type:       HINT_TYPE_TOKEN,
					TokenTypes: []int{tokens.TOKEN_WORD},
				},
			},
		},
		{
			Name: `io_here`,
			ParserHints: []*ParserHint{
				{
					Type:       HINT_TYPE_TOKEN,
					TokenTypes: []int{tokens.TOKEN_DLESS, tokens.TOKEN_DLESSDASH},
				},
				{
					Type:       HINT_TYPE_TOKEN,
					TokenTypes: []int{tokens.TOKEN_WORD},
				},
			},
		},
	})
}
