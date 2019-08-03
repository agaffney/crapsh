package lang

func init() {
	registerElements([]*ElementEntry{
		{
			Name: `SubshellCapture`,
			ParserData: []*ParserHint{
				{
					Type: HINT_TYPE_ANY,
					Members: []*ParserHint{
						{
							Type: HINT_TYPE_GROUP,
							Members: []*ParserHint{
								{
									Type:  HINT_TYPE_TOKEN,
									Name:  `SubshellCaptureOpen`,
									Final: true,
								},
								{
									Type: HINT_TYPE_ELEMENT,
									Name: `FullCommand`,
									Many: true,
								},
								{
									Type: HINT_TYPE_TOKEN,
									Name: `ParenClose`,
								},
							},
						},
						{
							Type: HINT_TYPE_GROUP,
							Members: []*ParserHint{
								{
									Type:  HINT_TYPE_TOKEN,
									Name:  `Backtick`,
									Final: true,
								},
								{
									Type: HINT_TYPE_ELEMENT,
									// TODO: Use variation of FullCommand that won't allow a backtick subshell
									Name: `FullCommand`,
									Many: true,
								},
								{
									Type: HINT_TYPE_TOKEN,
									Name: `Backtick`,
								},
							},
						},
					},
				},
			},
		},
		{
			Name: `Subshell`,
			ParserData: []*ParserHint{
				{
					Type: HINT_TYPE_GROUP,
					Members: []*ParserHint{
						{
							Type: HINT_TYPE_TOKEN,
							Name: `ParenOpen`,
						},
						{
							Type: HINT_TYPE_ELEMENT,
							Name: `FullCommand`,
							Many: true,
						},
						{
							Type: HINT_TYPE_TOKEN,
							Name: `ParenClose`,
						},
					},
				},
			},
		},
	})
	//registerParserHints([]*ParserHint{
	//	{
	//		Name:            `SubshellCapture`,
	//		TokenStart:      tokens.SUBSHELL_OPEN,
	//		TokenEnd:        tokens.PAREN_CLOSE,
	//		AllowedElements: []string{"Line"},
	//	},
	//	{
	//		Name:            `Subshell`,
	//		TokenStart:      tokens.PAREN_OPEN,
	//		TokenEnd:        tokens.PAREN_CLOSE,
	//		AllowedElements: []string{"Line"},
	//	},
	//	{
	//		Name:            `SubshellBacktick`,
	//		TokenStart:      tokens.BACKTICK,
	//		TokenEnd:        tokens.BACKTICK,
	//		AllowedElements: []string{"Line"},
	//	},
	//})
}
