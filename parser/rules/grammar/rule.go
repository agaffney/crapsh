package grammar

import (
	"github.com/agaffney/crapsh/parser/ast"
)

const (
	HINT_TYPE_NODE = iota
	HINT_TYPE_RULE
	HINT_TYPE_TOKEN
	HINT_TYPE_GROUP
	HINT_TYPE_ANY
)

type ParserHint struct {
	Type       int
	RuleName   string        // Name of rule or token to match
	Optional   bool          // Hint is optional
	Many       bool          // Hint can match multiple times
	Final      bool          // Consider the rule matched if this hint matches
	TokenTypes []int         // List of token types to match (for TOKEN type)
	Members    []*ParserHint // Child parser hints (used by ANY/GROUP hint types)
}

type GrammarRule struct {
	Name                   string          // name of the rule, used to refer to other rules from a parser hint
	ParserHints            []*ParserHint   // parser hints for the rule
	AstFunc                func() ast.Node `json:"-"` // don't include in JSON output, as it breaks encoding
	AllowFirstWordReserved bool            // whether the first word can be a reserved word
}

var GrammarRules = []*GrammarRule{}

func registerRules(rules []*GrammarRule) {
	GrammarRules = append(GrammarRules, rules...)
}

func GetRule(ruleName string) *GrammarRule {
	for _, rule := range GrammarRules {
		if rule.Name == ruleName {
			return rule
		}
	}
	return nil
}
