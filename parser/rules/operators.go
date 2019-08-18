package rules

import (
	"github.com/agaffney/crapsh/parser/tokens"
)

type OperatorRule struct {
	Pattern          string // string to match
	TokenType        int    // token type to return
	DelimitsIoNumber bool   // whether the operator delimits an IO_NUMBER token (starts with < or >)
}

// Helpful vim command to convert from parser/tokens format:
// %s/\(TOKEN_[^ ]\+\) \+\/\/ \(.*\)/{Pattern: `\2`, TokenType: tokens.\1},/
var OperatorRules = []OperatorRule{
	{Pattern: `;`, TokenType: tokens.TOKEN_SEMI},
	{Pattern: `|`, TokenType: tokens.TOKEN_PIPE},
	{Pattern: `&`, TokenType: tokens.TOKEN_AND},
	// TODO: maybe move these to delimeters
	// Disabled because they were breaking up the $( delimeter
	//	{Pattern: `(`, TokenType: tokens.TOKEN_LPAREN},
	//	{Pattern: `)`, TokenType: tokens.TOKEN_RPAREN},
	{Pattern: `<`, TokenType: tokens.TOKEN_LESS, DelimitsIoNumber: true},
	{Pattern: `>`, TokenType: tokens.TOKEN_GREAT, DelimitsIoNumber: true},
	{Pattern: `&&`, TokenType: tokens.TOKEN_AND_IF},
	{Pattern: `||`, TokenType: tokens.TOKEN_OR_IF},
	{Pattern: `;;`, TokenType: tokens.TOKEN_DSEMI},
	{Pattern: `<<`, TokenType: tokens.TOKEN_DLESS, DelimitsIoNumber: true},
	{Pattern: `>>`, TokenType: tokens.TOKEN_DGREAT, DelimitsIoNumber: true},
	{Pattern: `<&`, TokenType: tokens.TOKEN_LESSAND, DelimitsIoNumber: true},
	{Pattern: `>&`, TokenType: tokens.TOKEN_GREATAND, DelimitsIoNumber: true},
	{Pattern: `<>`, TokenType: tokens.TOKEN_LESSGREAT, DelimitsIoNumber: true},
	{Pattern: `<<-`, TokenType: tokens.TOKEN_DLESSDASH, DelimitsIoNumber: true},
	{Pattern: `>|`, TokenType: tokens.TOKEN_CLOBBER, DelimitsIoNumber: true},
}
