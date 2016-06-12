package tokens

import (
	"bytes"
	//"fmt"
	"regexp"
	"unicode"
)

const (
	DOLLAR_SIGN          = `$`
	PAREN_OPEN           = `(`
	PAREN_CLOSE          = `)`
	SQUARE_BRACKET_OPEN  = `[`
	SQUARE_BRACKET_CLOSE = `]`
	CURLY_BRACE_OPEN     = `{`
	CURLY_BRACE_CLOSE    = `}`
	BACKSLASH            = `\`
	SEMICOLON            = `;`
	BACKTICK             = "`"
	SUBSHELL_OPEN        = `$(`
	VARIABLE_OPEN        = `${`
	ARITHMETIC_OPEN      = `$((`
	DOUBLE_QUOTE         = `"`
	SINGLE_QUOTE         = "'"
	NEWLINE              = "\n"
	CARRIAGE_RETURN      = "\r"
	TAB                  = "\t"
	SPACE                = ` `
	KEYWORD_IF           = `if`
	KEYWORD_THEN         = `then`
	KEYWORD_ELSE         = `else`
	KEYWORD_FI           = `fi`
	KEYWORD_WHILE        = `while`
	KEYWORD_DO           = `do`
	KEYWORD_DONE         = `done`
	KEYWORD_FOR          = `for`
	KEYWORD_EVAL         = `eval`
)

const (
	TYPE_SIMPLE = iota
	TYPE_WHITESPACE
	TYPE_REGEXP
	TYPE_CALLBACK
	TYPE_MATCHALL
)

type Token struct {
	Type                int
	Name                string
	Pattern             string
	MatchUntilNextToken bool
	AdvanceLine         bool
}

func (t *Token) Match(buf *bytes.Buffer) (int, int) {
	switch {
	case t.Type == TYPE_SIMPLE:
		token_len := len(t.Pattern)
		if token_len == 0 {
			return -1, 0
		}
		buf_len := buf.Len()
		if buf_len < token_len {
			return -1, 0
		}
		//fmt.Printf("buf_len = %d, token_len = %d\n", buf_len, token_len)
		buf_bytes := buf.Bytes()[buf.Len()-token_len:]
		for i, b := range []byte(t.Pattern) {
			if buf_bytes[i] != b {
				return -1, 0
			}
		}
		return buf.Len() - token_len, token_len
	case t.Type == TYPE_REGEXP:
		foo := regexp.MustCompile(t.Pattern)
		match := foo.FindIndex(buf.Bytes())
		if match == nil {
			return -1, 0
		}
		return match[0], match[1] - match[0]
	case t.Type == TYPE_MATCHALL:
		return 0, buf.Len()
	case t.Type == TYPE_WHITESPACE:
		start_idx := -1
		length := 0
		for idx, c := range buf.String() {
			if unicode.IsSpace(c) {
				if start_idx == -1 {
					start_idx = idx
				}
				length++
			} else {
				if start_idx > -1 {
					break
				}
			}
		}
		//fmt.Printf("whitespace - start_idx=%d, length=%d, buf='%s'\n", start_idx, length, buf.Bytes())
		return start_idx, length
		//case TYPE_CALLBACK:

	}
	panic("Unknown token type on match!")
}

var Tokens []*Token

func registerTokens(tokens []*Token) {
	if Tokens == nil {
		Tokens = make([]*Token, 0)
	}
	for _, token := range tokens {
		Tokens = append(Tokens, token)
	}
}

func GetToken(name string) *Token {
	for _, t := range Tokens {
		if t.Name == name {
			return t
		}
	}
	return nil
}

func init() {
	registerTokens([]*Token{
		{
			Name:    `Escape`,
			Type:    TYPE_REGEXP,
			Pattern: `\\.`,
		},
		{
			Name:        `Newline`,
			Pattern:     "\n",
			AdvanceLine: true,
		},
		{
			Name:                `Whitespace`,
			Type:                TYPE_WHITESPACE,
			MatchUntilNextToken: true,
		},
		{
			Name:    `Dollar`,
			Pattern: `$`,
		},
		{
			Name:    `SingleQuote`,
			Pattern: `'`,
		},
		{
			Name:    `DoubleQuote`,
			Pattern: `"`,
		},
		{
			Name:    `ParenOpen`,
			Pattern: `(`,
		},
		{
			Name:    `ParenClose`,
			Pattern: `)`,
		},
		{
			Name:    `CurlyBraceOpen`,
			Pattern: `{`,
		},
		{
			Name:    `CurlyBraceClose`,
			Pattern: `}`,
		},
		{
			Name:    `SquareBracketOpen`,
			Pattern: `[`,
		},
		{
			Name:    `SquareBracketClose`,
			Pattern: `]`,
		},
	})
}
