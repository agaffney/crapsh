package lang

func init() {
	registerElements([]*ElementEntry{
		{
			Name: `Root`,
			ParserData: []*ParserHint{
				{
					Type: HINT_TYPE_ELEMENT,
					Name: `Command`,
					Many: true,
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
				{
					Type:     HINT_TYPE_ANY,
					Optional: true,
					Members: []*ParserHint{
						{
							Type: HINT_TYPE_TOKEN,
							Name: `Newline`,
						},
						{
							Type: HINT_TYPE_TOKEN,
							Name: `Semicolon`,
						},
					},
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
						/*
							{
								Type: HINT_TYPE_ELEMENT,
								Name: `StringSingle`,
							},
							{
								Type: HINT_TYPE_ELEMENT,
								Name: `StringDouble`,
							},
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
					Name:     `Whitespace`,
					Optional: true,
				},
			},
		},
		{
			Name: `Generic`,
			ParserData: []*ParserHint{
				{
					Type: HINT_TYPE_TOKEN,
					//Name: `Generic`,
				},
			},
		},
	})

}
