package parser

import (
	"github.com/agaffney/crapsh/parser/lexer"
	"github.com/agaffney/crapsh/parser/rules"
	"github.com/agaffney/crapsh/parser/rules/grammar"
	"github.com/agaffney/crapsh/parser/tokens"
	"regexp"
)

// Check if existing token is a reversed word
func (p *Parser) checkTokenIsReserved(tokenType int) *rules.ReservedRule {
	for _, rule := range rules.ReservedRules {
		if tokenType == rule.TokenType {
			return &rule
		}
	}
	return nil
}

// Classify token based on current parser hint
func (p *Parser) classifyToken(token *lexer.Token, hint *grammar.ParserHint) int {
	// No need to classify if we already have
	if token.Type != tokens.TOKEN_NULL {
		return token.Type
	}
	if p.stack.Cur().allowNextWordReserved {
		p.stack.Cur().allowNextWordReserved = false
		for _, rule := range rules.ReservedRules {
			// TODO: check for:
			// * previous tokens match AfterTokens (-1 is wildcard)
			if token.Value == rule.Pattern {
				return rule.TokenType
			}
		}
	}
	NAME_RE := regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]+$`)
	ASSIGNMENT_RE := regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]+=`)
	for hintTokenType := range hint.TokenTypes {
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
