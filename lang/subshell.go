package lang

func init() {
	registerElements([]*ElementEntry{
		{
			Name: `SubshellCaptureParens`,
			ParserData: []*ParserHint{
				{
					Type:   HINT_TYPE_TOKEN,
					Tokens: []string{`SubshellCaptureOpen`},
					Final:  true,
				},
				{
					Type: HINT_TYPE_ELEMENT,
					Name: `FullCommand`,
					Many: true,
				},
				{
					Type:   HINT_TYPE_TOKEN,
					Tokens: []string{`ParenClose`},
				},
			},
		},
		{
			Name: `SubshellCaptureBackticks`,
			ParserData: []*ParserHint{
				{
					Type:   HINT_TYPE_TOKEN,
					Tokens: []string{`Backtick`},
					Final:  true,
				},
				{
					Type: HINT_TYPE_ELEMENT,
					// TODO: Use variation of FullCommand that won't allow a backtick subshell
					Name: `FullCommand`,
					Many: true,
				},
				{
					Type:   HINT_TYPE_TOKEN,
					Tokens: []string{`Backtick`},
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
							Type:   HINT_TYPE_TOKEN,
							Tokens: []string{`ParenOpen`},
						},
						{
							Type: HINT_TYPE_ELEMENT,
							Name: `FullCommand`,
							Many: true,
						},
						{
							Type:   HINT_TYPE_TOKEN,
							Tokens: []string{`ParenClose`},
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
