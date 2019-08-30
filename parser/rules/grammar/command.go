package grammar

import (
	"github.com/agaffney/crapsh/parser/ast"
	"github.com/agaffney/crapsh/parser/tokens"
)

func init() {
	registerRules([]*GrammarRule{
		{
			Name:    `complete_command`,
			AstFunc: ast.NewCompleteCommand,
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
		{
			Name:    `pipeline`,
			AstFunc: ast.NewPipeline,
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
			Name:    `simple_command`,
			AstFunc: ast.NewSimpleCommand,
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
								},
								{
									Type:     HINT_TYPE_GROUP,
									Optional: true,
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
					Many: true,
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
