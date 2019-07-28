package lexer

import ()

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

// Attempt to match token definition against buffer
// Returns index and length
func (t *TokenDefinition) Match(buf *bytes.Buffer) (int, int) {
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

var Tokens = []TokenDefinition{
	{
		Name:    `DollarSign`,
		Pattern: `$`,
	},
}
