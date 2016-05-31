package parser

import (
	"fmt"
	"github.com/agaffney/crapsh/lang/tokens"
	"io"
)

type Token struct {
	Type  string
	Value string
}

func (p *Parser) nextToken() (*Token, error) {
	for {
		// Check the buffer at the beginning to catch tokens already in the buffer
		// from the last iteration
		if p.buf.Len() > 0 {
			for _, foo := range tokens.Tokens {
				if idx := foo.Match(p.buf.Buffer); idx > -1 {
					var token *Token
					if idx == 0 {
						token = &Token{Type: foo.Name, Value: p.buf.String()}
						p.buf.Reset()
					} else {
						// Return data up to token as "generic" token and remove from buffer
						token = &Token{Type: `Generic`, Value: string(p.buf.Bytes()[0 : idx-1])}
						p.buf = NewBuffer(p.buf.Bytes()[idx:])
					}
					return token, nil
				}
			}
		}
		c, err := p.nextRune()
		fmt.Printf("Line %d, offset %d: %#U\n", p.Position.Line, p.Position.LineOffset, c)
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			break
		}
		p.buf.WriteRune(c)
	}
	return nil, nil
}
