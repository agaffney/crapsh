package grammar

import (
	"github.com/agaffney/crapsh/parser/tokens"
)

func init() {
	registerRules([]*GrammarRule{
		{
			Name: `newline_list`,
			ParserHints: []*ParserHint{
				{
					Type:       HINT_TYPE_TOKEN,
					TokenTypes: []int{tokens.TOKEN_NEWLINE},
					Many:       true,
				},
			},
		},
		{
			Name: `linebreak`,
			ParserHints: []*ParserHint{
				{
					Type:     HINT_TYPE_RULE,
					RuleName: `newline_list`,
					Optional: true,
				},
			},
		},
		{
			Name: `separator_op`,
			ParserHints: []*ParserHint{
				{
					Type:       HINT_TYPE_TOKEN,
					TokenTypes: []int{tokens.TOKEN_AND, tokens.TOKEN_SEMI},
				},
			},
		},
		{
			Name: `separator`,
			ParserHints: []*ParserHint{
				{
					Type: HINT_TYPE_ANY,
					Members: []*ParserHint{
						{
							Type: HINT_TYPE_GROUP,
							Members: []*ParserHint{
								{
									Type:     HINT_TYPE_RULE,
									RuleName: `separator_op`,
								},
								{
									Type:     HINT_TYPE_RULE,
									RuleName: `linebreak`,
								},
							},
						},
						{
							Type:     HINT_TYPE_RULE,
							RuleName: `newline_list`,
						},
					},
				},
			},
		},
		{
			Name: `sequential_sep`,
			ParserHints: []*ParserHint{
				{
					Type: HINT_TYPE_ANY,
					Members: []*ParserHint{
						{
							Type: HINT_TYPE_GROUP,
							Members: []*ParserHint{
								{
									Type:       HINT_TYPE_TOKEN,
									TokenTypes: []int{tokens.TOKEN_SEMI},
								},
								{
									Type:     HINT_TYPE_RULE,
									RuleName: `linebreak`,
								},
							},
						},
						{
							Type:     HINT_TYPE_RULE,
							RuleName: `newline_list`,
						},
					},
				},
			},
		},
	})
}
