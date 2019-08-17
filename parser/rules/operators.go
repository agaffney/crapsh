package rules

import (
	"github.com/agaffney/crapsh/parser/tokens"
)

type OperatorRule struct {
	Pattern   string
	TokenType int
}

// Helpful vim command to convert from parser/tokens format:
// %s/\(TOKEN_[^ ]\+\) \+\/\/ \(.*\)/{Pattern: `\2`, TokenType: tokens.\1},/
var OperatorRules = []OperatorRule{
	{Pattern: `&&`, TokenType: tokens.TOKEN_AND_IF},
	{Pattern: `||`, TokenType: tokens.TOKEN_OR_IF},
	{Pattern: `;;`, TokenType: tokens.TOKEN_DSEMI},
	{Pattern: `<<`, TokenType: tokens.TOKEN_DLESS},
	{Pattern: `>>`, TokenType: tokens.TOKEN_DGREAT},
	{Pattern: `<&`, TokenType: tokens.TOKEN_LESSAND},
	{Pattern: `>&`, TokenType: tokens.TOKEN_GREATAND},
	{Pattern: `<>`, TokenType: tokens.TOKEN_LESSGREAT},
	{Pattern: `<<-`, TokenType: tokens.TOKEN_DLESSDASH},
	{Pattern: `>|`, TokenType: tokens.TOKEN_CLOBBER},
	{Pattern: `;`, TokenType: tokens.TOKEN_SEMI},
	{Pattern: `|`, TokenType: tokens.TOKEN_PIPE},
	{Pattern: `<`, TokenType: tokens.TOKEN_LESS},
	{Pattern: `>`, TokenType: tokens.TOKEN_GREAT},
	{Pattern: `&`, TokenType: tokens.TOKEN_AND},
	// TODO: maybe move these to delimeters
	// Disabled because they were breaking up the $( delimeter
	//	{Pattern: `(`, TokenType: tokens.TOKEN_LPAREN},
	//	{Pattern: `)`, TokenType: tokens.TOKEN_RPAREN},
}
