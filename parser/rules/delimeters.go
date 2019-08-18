package rules

type DelimeterRule struct {
	Name           string   // rule name
	DelimStart     string   // starting delimeter
	DelimEnd       string   // ending delimeter
	AllowedRules   []string // delimeter rule names that are allowed inside this delimeter set
	ReturnToken    bool     // whether to return a token when finding the end delimeter
	IgnoreEscapes  bool     // whether to ignore escape sequences
	IncludeDelim   bool     // whether to include the delimeters in the token
	AllowOperators bool     // whether to allow operators
}

var DelimeterRules = []DelimeterRule{
	{
		Name:           `Word`,
		DelimEnd:       ` `,
		ReturnToken:    true,
		AllowOperators: true,
		AllowedRules:   []string{`SingleQuotes`, `DoubleQuotes`, `Subshell`, `SubshellBackticks`, `Arithmetic`},
	},
	{
		Name:          `SingleQuotes`,
		DelimStart:    `'`,
		DelimEnd:      `'`,
		IgnoreEscapes: true,
		IncludeDelim:  true,
	},
	{
		Name:         `DoubleQuotes`,
		DelimStart:   `"`,
		DelimEnd:     `"`,
		IncludeDelim: true,
		AllowedRules: []string{`Subshell`, `SubshellBackticks`, `Arithmetic`},
	},
	{
		Name:         `Subshell`,
		DelimStart:   `$(`,
		DelimEnd:     `)`,
		IncludeDelim: true,
		AllowedRules: []string{`SingleQuotes`, `DoubleQuotes`, `Subshell`, `SubshellBackticks`, `Arithmetic`},
	},
	{
		Name:         `SubshellBackticks`,
		DelimStart:   "`",
		DelimEnd:     "`",
		IncludeDelim: true,
		AllowedRules: []string{`SingleQuotes`, `DoubleQuotes`, `Subshell`, `Arithmetic`},
	},
	{
		Name:         `Arithmetic`,
		DelimStart:   `$((`,
		DelimEnd:     `))`,
		IncludeDelim: true,
	},
}

func GetDelimeterRule(name string) *DelimeterRule {
	for _, rule := range DelimeterRules {
		if name == rule.Name {
			return &rule
		}
	}
	return nil
}
