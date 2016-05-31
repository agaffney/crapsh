package tokens

import (
	"bytes"
	"regexp"
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
	TYPE_SIMPLE = 0
	TYPE_REGEXP
	TYPE_CALLBACK
	TYPE_MATCHALL
)

type Token struct {
	Type    int
	Name    string
	Pattern string
}

func (t *Token) Match(buf *bytes.Buffer) int {
	switch {
	case t.Type == TYPE_SIMPLE:
		token_len := len(t.Pattern)
		if token_len == 0 {
			return -1
		}
		buf_len := buf.Len()
		if buf_len < token_len {
			return -1
		}
		//fmt.Printf("buf_len = %d, token_len = %d\n", buf_len, token_len)
		buf_bytes := buf.Bytes()[buf.Len()-token_len:]
		for i, b := range []byte(t.Pattern) {
			if buf_bytes[i] != b {
				return -1
			}
		}
		return buf.Len() - token_len
	case t.Type == TYPE_REGEXP:
		foo := regexp.MustCompile(t.Pattern)
		match := foo.FindIndex(buf.Bytes())
		return match[0]
	case t.Type == TYPE_MATCHALL:
		return 0
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

func init() {
	registerTokens([]*Token{
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
			Name: `Generic`,
			Type: TYPE_MATCHALL,
		},
	})
}
