package rules

import (
	"github.com/agaffney/crapsh/parser/tokens"
)

type ReservedRule struct {
	Pattern   string // string to look for
	TokenType int    // token type to return
}

var ReservedRules = []ReservedRule{
	{Pattern: `if`, TokenType: tokens.TOKEN_IF},
	{Pattern: `then`, TokenType: tokens.TOKEN_THEN},
	{Pattern: `else`, TokenType: tokens.TOKEN_ELSE},
	{Pattern: `elif`, TokenType: tokens.TOKEN_ELIF},
	{Pattern: `fi`, TokenType: tokens.TOKEN_FI},
	{Pattern: `do`, TokenType: tokens.TOKEN_DO},
	{Pattern: `done`, TokenType: tokens.TOKEN_DONE},
	{Pattern: `case`, TokenType: tokens.TOKEN_CASE},
	{Pattern: `esac`, TokenType: tokens.TOKEN_ESAC},
	{Pattern: `while`, TokenType: tokens.TOKEN_WHILE},
	{Pattern: `until`, TokenType: tokens.TOKEN_UNTIL},
	{Pattern: `for`, TokenType: tokens.TOKEN_FOR},
	{Pattern: `{`, TokenType: tokens.TOKEN_LBRACE},
	{Pattern: `}`, TokenType: tokens.TOKEN_RBRACE},
	{Pattern: `!`, TokenType: tokens.TOKEN_BANG},
	{Pattern: `in`, TokenType: tokens.TOKEN_IN},
}
