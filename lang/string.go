package lang

const (
	STRING_TYPE_INVALID int = iota
	STRING_TYPE_SINGLE
	STRING_TYPE_DOUBLE
)

type String struct {
	*Generic
	Type int
}

func NewStringSingle(base *Generic) Element {
	return &String{Generic: base, Type: STRING_TYPE_SINGLE}
}

func NewStringDouble(base *Generic) Element {
	return &String{Generic: base, Type: STRING_TYPE_DOUBLE}
}

func init() {
	registerElements([]*ElementEntry{
		{
			Name: `StringSingle`,
			ParserData: []*ParserHint{
				{
					Type:   HINT_TYPE_TOKEN,
					Tokens: []string{`SingleQuote`},
					Final:  true,
				},
				// TODO: list of valid tokens
				{
					Type: HINT_TYPE_TOKEN,
					//Name: `Generic`,
					Many: true,
				},
				{
					Type:   HINT_TYPE_TOKEN,
					Tokens: []string{`SingleQuote`},
				},
			},
			Factory: NewStringSingle,
		},
		{
			Name: `StringDouble`,
			ParserData: []*ParserHint{
				{
					Type:   HINT_TYPE_TOKEN,
					Tokens: []string{`DoubleQuote`},
					Final:  true,
				},
				// TODO: list of valid tokens
				{
					Type: HINT_TYPE_TOKEN,
					//Name: `Generic`,
					Many: true,
				},
				{
					Type:   HINT_TYPE_TOKEN,
					Tokens: []string{`DoubleQuote`},
				},
			},
			//AllowedElements: []string{"Variable", "SubshellCapture", "Subshell", "SubshellBacktick", "Generic"},
			Factory: NewStringDouble,
		},
	})
}
