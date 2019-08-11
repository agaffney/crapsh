package rules

import (
	"github.com/agaffney/crapsh/parser/tokens"
)

type reservedRule struct {
	tokenType int
	pattern   string
}

var reservedRules = []reservedRule{
	{pattern: `if`, tokenType: tokens.TOKEN_IF},
	{pattern: `then`, tokenType: tokens.TOKEN_THEN},
	{pattern: `else`, tokenType: tokens.TOKEN_ELSE},
	{pattern: `elif`, tokenType: tokens.TOKEN_ELIF},
	{pattern: `fi`, tokenType: tokens.TOKEN_FI},
	{pattern: `do`, tokenType: tokens.TOKEN_DO},
	{pattern: `done`, tokenType: tokens.TOKEN_DONE},
	{pattern: `case`, tokenType: tokens.TOKEN_CASE},
	{pattern: `esac`, tokenType: tokens.TOKEN_ESAC},
	{pattern: `while`, tokenType: tokens.TOKEN_WHILE},
	{pattern: `until`, tokenType: tokens.TOKEN_UNTIL},
	{pattern: `for`, tokenType: tokens.TOKEN_FOR},
	{pattern: `{`, tokenType: tokens.TOKEN_LBRACE},
	{pattern: `}`, tokenType: tokens.TOKEN_RBRACE},
	{pattern: `!`, tokenType: tokens.TOKEN_BANG},
	{pattern: `in`, tokenType: tokens.TOKEN_IN},
}
