package parser

import (
	"github.com/agaffney/crapsh/parser/lexer"
	"github.com/agaffney/crapsh/parser/rules"
	"github.com/agaffney/crapsh/parser/rules/grammar"
	"github.com/agaffney/crapsh/parser/tokens"
	"regexp"
)

// Classify token based on current parser hint
func (p *Parser) classifyToken(token *lexer.Token, hint *grammar.ParserHint) int {
	// No need to classify if we already have
	if token.Type != tokens.TOKEN_NULL {
		return token.Type
	}
	NAME_RE := regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]+$`)
	ASSIGNMENT_RE := regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]+=`)
	for hintTokenType := range hint.TokenTypes {
		for _, rule := range rules.ReservedRules {
			if token.Value == rule.Pattern {
				return rule.TokenType
			}
		}
		switch hintTokenType {
		case tokens.TOKEN_NAME:
			if NAME_RE.MatchString(token.Value) {
				return tokens.TOKEN_NAME
			}
		case tokens.TOKEN_ASSIGNMENT_WORD:
			if ASSIGNMENT_RE.MatchString(token.Value) {
				return tokens.TOKEN_ASSIGNMENT_WORD
			}
		}
	}
	// Fallback to generic WORD token type
	return tokens.TOKEN_WORD
}
