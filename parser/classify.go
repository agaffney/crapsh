package parser

import (
	//"fmt"
	"github.com/agaffney/crapsh/parser/lexer"
	"github.com/agaffney/crapsh/parser/rules"
	"github.com/agaffney/crapsh/parser/rules/grammar"
	"github.com/agaffney/crapsh/parser/tokens"
	//"github.com/agaffney/crapsh/util"
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
	//util.DumpObject(hint, "classifyToken(): hint = ")
	// No need to classify if we already have
	if token.Type != tokens.TOKEN_NULL {
		//fmt.Printf("classifyToken(): returning existing token type\n")
		return token.Type
	}
	if p.stack.Cur().allowNextWordReserved {
		p.stack.Cur().allowNextWordReserved = false
		for _, rule := range rules.ReservedRules {
			// TODO: check for:
			// * previous tokens match AfterTokens (-1 is wildcard), for in/do after case/for
			if token.Value == rule.Pattern {
				return rule.TokenType
			}
		}
	}
	NAME_RE := regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]+$`)
	ASSIGNMENT_RE := regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]+=`)
	for _, hintTokenType := range hint.TokenTypes {
		switch hintTokenType {
		case tokens.TOKEN_NAME:
			if NAME_RE.MatchString(token.Value) {
				return tokens.TOKEN_NAME
			}
		case tokens.TOKEN_ASSIGNMENT_WORD:
			//fmt.Printf("classifyToken(): checking '%s' for ASSIGNMENT_WORD\n", token.Value)
			if ASSIGNMENT_RE.MatchString(token.Value) {
				return tokens.TOKEN_ASSIGNMENT_WORD
			}
		}
	}
	// Fallback to generic WORD token type
	return tokens.TOKEN_WORD
}
