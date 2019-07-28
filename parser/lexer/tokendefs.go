package lexer

import (
	"bytes"
	"regexp"
	"unicode"
)

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

func (t *TokenDefinition) findNextToken(buf *bytes.Buffer, offset int) int {
	for i := offset; i < len(buf.String()); i++ {
		offset++
		for _, foo := range TokenDefinitions {
			if foo.Name == t.Name {
				continue
			}
			// This will match anything, so we skip when looking for the next token
			if foo.Type == TYPE_MATCHALL {
				continue
			}
			if ok, _ := foo.Match(buf, offset); ok {
				return offset
			}
		}
	}
	return -1
}

// Attempt to match token definition against buffer
func (t *TokenDefinition) Match(buf *bytes.Buffer, offset int) (bool, string) {
	switch {
	case t.Type == TYPE_SIMPLE:
		token_len := len(t.Pattern)
		if token_len == 0 {
			return false, ""
		}
		buf_len := len(buf.String())
		if buf_len < (offset + token_len) {
			return false, ""
		}
		buf_str := buf.String()[offset : offset+token_len]
		if buf_str != t.Pattern {
			return false, ""
		}
		return true, buf_str
	case t.Type == TYPE_REGEXP:
		foo := regexp.MustCompile(t.Pattern)
		match := foo.FindStringIndex(buf.String())
		if match == nil {
			return false, ""
		}
		return true, buf.String()[match[0] : match[1]-match[0]]
	case t.Type == TYPE_MATCHALL:
		if t.MatchUntilNextToken {
			nextOffset := t.findNextToken(buf, offset)
			if nextOffset > 0 {
				return true, buf.String()[offset : nextOffset-offset]
			} else {
				return true, buf.String()[offset:]
			}
		} else {
			return true, buf.String()[offset:]
		}
	case t.Type == TYPE_WHITESPACE:
		ret := bytes.NewBuffer(nil)
		for _, c := range buf.String()[offset:] {
			if unicode.IsSpace(c) {
				ret.WriteRune(c)
			} else {
				break
			}
		}
		if ret.Len() > 0 {
			return true, ret.String()
		} else {
			return false, ""
		}
	case t.Type == TYPE_CALLBACK:
		return false, ""
	}
	panic("Unknown token type on match!")
}

var TokenDefinitions = []TokenDefinition{
	{
		Name:        `Newline`,
		Pattern:     "\n",
		AdvanceLine: true,
	},
	{
		Name:    `DollarSign`,
		Pattern: `$`,
	},
	{
		Name:    `Semicolon`,
		Pattern: `;`,
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
		Name:    `BackTick`,
		Pattern: "`",
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
	{
		Name:                `Whitespace`,
		Type:                TYPE_WHITESPACE,
		MatchUntilNextToken: true,
	},
	{
		Name:        `EscapeNewline`,
		Pattern:     "\\\n",
		AdvanceLine: true,
	},
	{
		Name: `Escape`,
		Type: TYPE_REGEXP,
		// Backslash followed by any character
		Pattern: `\\.`,
	},
	{
		Name:    `Identifier`,
		Type:    TYPE_REGEXP,
		Pattern: `[a-zA-Z_][a-zA-Z0-9_]+`,
	},
	{
		Name:                `Generic`,
		Type:                TYPE_MATCHALL,
		MatchUntilNextToken: true,
	},
}
