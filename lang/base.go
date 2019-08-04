package lang

func init() {
	registerElements([]*ElementEntry{
		{
			Name: `Root`,
			ParserData: []*ParserHint{
				{
					Type: HINT_TYPE_ELEMENT,
					Name: `Command`,
				},
				{
					Type:     HINT_TYPE_GROUP,
					Many:     true,
					Optional: true,
					Members: []*ParserHint{
						{
							Type:   HINT_TYPE_TOKEN,
							Tokens: []string{`Pipe`},
							Final:  true,
						},
						{
							Type: HINT_TYPE_ELEMENT,
							Name: `Command`,
						},
					},
				},
			},
		},
		{
			Name: `Command`,
			ParserData: []*ParserHint{
				{
					Type: HINT_TYPE_ANY,
					Members: []*ParserHint{
						{
							Type: HINT_TYPE_ELEMENT,
							Name: `Subshell`,
						},
						{
							Type: HINT_TYPE_ELEMENT,
							Name: `Argument`,
							Many: true,
						},
					},
				},
				// TODO: move to FullCommand
				{
					Type:     HINT_TYPE_TOKEN,
					Tokens:   []string{`Newline`, `Semicolon`},
					Optional: true,
				},
			},
		},
		{
			Name: `Argument`,
			ParserData: []*ParserHint{
				{
					Type: HINT_TYPE_ANY,
					Many: true,
					Members: []*ParserHint{
						{
							Type: HINT_TYPE_ELEMENT,
							Name: `StringSingle`,
						},
						{
							Type: HINT_TYPE_ELEMENT,
							Name: `StringDouble`,
						},
						/*
							{
								Type: HINT_TYPE_ELEMENT,
								Name: `SubshellCapture`,
							},
						*/
						{
							Type: HINT_TYPE_ELEMENT,
							Name: `Generic`,
						},
					},
				},
				{
					Type:     HINT_TYPE_TOKEN,
					Tokens:   []string{`Whitespace`},
					Optional: true,
				},
			},
		},
		// TODO: specify list of tokens
		{
			Name: `Generic`,
			ParserData: []*ParserHint{
				{
					Type: HINT_TYPE_TOKEN,
					//Name: `Generic`,
					Tokens: []string{`Generic`, `GenericIdentifier`},
				},
			},
		},
	})

}
