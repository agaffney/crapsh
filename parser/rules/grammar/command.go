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
		// 'complete_command' and 'list' have been merged together
		{
			Name: `complete_command`,
			ParserHints: []*ParserHint{
				{
					Type: HINT_TYPE_GROUP,
					Many: true,
					Members: []*ParserHint{
						{
							Type:     HINT_TYPE_RULE,
							RuleName: `and_or`,
						},
						{
							Type:     HINT_TYPE_RULE,
							RuleName: `separator_op`,
							Optional: true,
						},
					},
				},
			},
		},
		{
			Name: `and_or`,
			ParserHints: []*ParserHint{
				{
					Type:     HINT_TYPE_RULE,
					RuleName: `pipeline`,
				},
				{
					Type:     HINT_TYPE_GROUP,
					Optional: true,
					Many:     true,
					Members: []*ParserHint{
						{
							Type:       HINT_TYPE_TOKEN,
							TokenTypes: []int{tokens.TOKEN_AND_IF, tokens.TOKEN_OR_IF},
						},
						{
							Type:     HINT_TYPE_RULE,
							RuleName: `linebreak`,
						},
						{
							Type:     HINT_TYPE_RULE,
							RuleName: `pipeline`,
						},
					},
				},
			},
		},
		// 'pipeline' and 'pipe_sequence' have been merged together
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
									RuleName: `compound_command`,
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
		{
			Name: `simple_command`,
			ParserHints: []*ParserHint{
				{
					Type: HINT_TYPE_ANY,
					Members: []*ParserHint{
						{
							Type: HINT_TYPE_GROUP,
							Members: []*ParserHint{
								{
									Type:     HINT_TYPE_RULE,
									RuleName: `cmd_prefix`,
									Many:     true,
								},
								{
									Type:     HINT_TYPE_GROUP,
									Optional: true,
									Members: []*ParserHint{
										{
											Type:     HINT_TYPE_RULE,
											RuleName: `cmd_word`,
										},
										{
											Type:     HINT_TYPE_RULE,
											RuleName: `cmd_suffix`,
											Optional: true,
										},
									},
								},
							},
						},
						{
							Type: HINT_TYPE_GROUP,
							Members: []*ParserHint{
								{
									Type:       HINT_TYPE_TOKEN,
									TokenTypes: []int{tokens.TOKEN_WORD},
								},
								{
									Type:     HINT_TYPE_RULE,
									RuleName: `cmd_suffix`,
									Optional: true,
								},
							},
						},
					},
				},
			},
		},
		{
			Name: `cmd_prefix`,
			ParserHints: []*ParserHint{
				{
					Type: HINT_TYPE_ANY,
					Members: []*ParserHint{
						{
							Type:       HINT_TYPE_TOKEN,
							TokenTypes: []int{tokens.TOKEN_ASSIGNMENT_WORD},
						},
						{
							Type:     HINT_TYPE_RULE,
							RuleName: `io_redirect`,
						},
					},
				},
			},
		},
		{
			Name: `cmd_suffix`,
			ParserHints: []*ParserHint{
				{
					Type: HINT_TYPE_ANY,
					Many: true,
					Members: []*ParserHint{
						{
							Type:       HINT_TYPE_TOKEN,
							TokenTypes: []int{tokens.TOKEN_WORD},
						},
						{
							Type:     HINT_TYPE_RULE,
							RuleName: `io_redirect`,
						},
					},
				},
			},
		},
	})
}