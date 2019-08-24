package grammar

import (
	"github.com/agaffney/crapsh/parser/tokens"
)

func init() {
	registerRules([]*GrammarRule{
		{
			Name:                   `BasicCommand`,
			AllowFirstWordReserved: true,
			ParserHints: []*ParserHint{
				{
					Type:       HINT_TYPE_TOKEN,
					TokenTypes: []int{tokens.TOKEN_WORD},
					Many:       true,
				},
				{
					Type:       HINT_TYPE_TOKEN,
					TokenTypes: []int{tokens.TOKEN_NEWLINE},
					Optional:   true,
				},
			},
		},
		// TODO: merge 'list' and 'and_or' rules into this one
		{
			Name: `complete_command`,
			ParserHints: []*ParserHint{
				{
					Type:     HINT_TYPE_RULE,
					RuleName: `list`,
				},
				{
					Type:     HINT_TYPE_RULE,
					RuleName: `separator_op`,
					Optional: true,
				},
			},
		},
		// pipeline and pipe_sequence have been merged together
		{
			Name: `pipeline`,
			ParserHints: []*ParserHint{
				{
					Type:       HINT_TYPE_TOKEN,
					TokenTypes: []int{tokens.TOKEN_BANG},
					Optional:   true,
				},
				{
					Type:     HINT_TYPE_RULE,
					RuleName: `command`,
					Final:    true,
				},
				{
					Type:     HINT_TYPE_GROUP,
					Optional: true,
					Members: []*ParserHint{
						{
							Type:       HINT_TYPE_TOKEN,
							TokenTypes: []int{tokens.TOKEN_PIPE},
						},
						{
							Type:     HINT_TYPE_RULE,
							RuleName: `linebreak`,
						},
						{
							Type:     HINT_TYPE_RULE,
							RuleName: `command`,
						},
					},
				},
			},
		},
		{
			Name: `command`,
			ParserHints: []*ParserHint{
				{
					Type: HINT_TYPE_ANY,
					Members: []*ParserHint{
						{
							Type:     HINT_TYPE_RULE,
							RuleName: `simple_command`,
						},
						{
							Type:     HINT_TYPE_RULE,
							RuleName: `compound_command`,
						},
						{
							Type: HINT_TYPE_GROUP,
							Members: []*ParserHint{
								{
									Type:     HINT_TYPE_RULE,
									RuleName: `simple_command`,
								},
								{
									Type:     HINT_TYPE_RULE,
									RuleName: `redirect_list`,
								},
							},
						},
						{
							Type:     HINT_TYPE_RULE,
							RuleName: `function_definition`,
						},
					},
				},
			},
		},
	})
}
