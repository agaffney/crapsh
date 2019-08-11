package rules

type DelimeterRule struct {
	Name           string
	DelimStart     string
	DelimEnd       string
	AllowedRules   []string
	ReturnToken    bool
	IgnoreEscapes  bool
	IncludeDelim   bool
	AllowOperators bool
}

var DelimeterRules = []DelimeterRule{
	{Name: `Word`, DelimEnd: ` `, ReturnToken: true, AllowOperators: true, AllowedRules: []string{`SingleQuotes`, `DoubleQuotes`, `Subshell`, `SubshellBackticks`}},
	{Name: `SingleQuotes`, DelimStart: `'`, DelimEnd: `'`, IgnoreEscapes: true, IncludeDelim: true},
	{Name: `DoubleQuotes`, DelimStart: `"`, DelimEnd: `"`, IncludeDelim: true, AllowedRules: []string{}},
	{Name: `Subshell`, DelimStart: `$(`, DelimEnd: `)`, IncludeDelim: true, AllowedRules: []string{}},
	{Name: `SubshellBackticks`, DelimStart: "`", DelimEnd: "`", IncludeDelim: true, AllowedRules: []string{}},
}

func GetDelimeterRule(name string) *DelimeterRule {
	for _, rule := range DelimeterRules {
		if name == rule.Name {
			return &rule
		}
	}
	return nil
}
