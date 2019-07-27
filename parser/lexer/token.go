package lexer

import ()

type Token struct {
	LineNum int
	Offset  int
	Value   string
}

const (
	TYPE_SIMPLE = iota
	TYPE_WHITESPACE
	TYPE_REGEXP
	TYPE_CALLBACK
	TYPE_MATCHALL
)

type TokenDefinition struct {
	Type                int
	Name                string
	Pattern             string
	MatchUntilNextToken bool
	AdvanceLine         bool
}
