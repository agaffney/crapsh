package parser

import (
	"fmt"
	"github.com/agaffney/crapsh/lang/tokens"
	//"github.com/agaffney/crapsh/util"
	//"io"
)

type Token struct {
	Type  string
	Value string
	Pos   Position
}

func (p *Parser) getTokenIdx() int {
	return p.tokenIdx
}

func (p *Parser) setTokenIdx(idx int) {
	p.tokenIdx = idx
	if idx >= 0 {
		fmt.Printf("setTokenIdx(%d): curToken = %#v\n", idx, p.curToken())
	}
}

func (p *Parser) curToken() *Token {
	return p.tokenBuf[p.tokenIdx]
}

func (p *Parser) prevToken() {
	if p.tokenIdx > 0 {
		p.tokenIdx--
	}
}

func (p *Parser) nextToken() error {
	if p.tokenIdx < len(p.tokenBuf)-1 {
		p.tokenIdx++
	} else {
		token, err := p.readToken()
		if err != nil {
			return err
		}
		if token == nil {
			return nil
		}
		p.tokenBuf = append(p.tokenBuf, token)
		p.tokenIdx++
	}
	return nil
}

func (p *Parser) readToken() (*Token, error) {
	var buf_len int
	for {
		// Check the buffer at the beginning to catch tokens already in the buffer
		// from the last iteration
		buf_len = p.buf.Len()
		if buf_len > 0 {
			for _, foo := range tokens.Tokens {
				if idx, length := foo.Match(p.buf.Buffer); idx > -1 {
					//fmt.Printf("idx=%d, length=%d, data='%s'\n", idx, length, p.buf.Bytes()[idx:idx+length])
					if length == buf_len && foo.MatchUntilNextToken {
						break
					}
					var token *Token
					if idx > 0 {
						// Return data up to token as "generic" token and remove from buffer
						token = &Token{Type: `Generic`, Value: string(p.buf.Bytes()[0:idx]), Pos: p.Position}
						p.buf = NewBuffer(p.buf.Bytes()[idx:])
					} else {
						token = &Token{Type: foo.Name, Value: string(p.buf.Bytes()[idx : idx+length]), Pos: p.Position}
						if length == p.buf.Len() {
							p.buf.Reset()
						} else {
							p.buf = NewBuffer(p.buf.Bytes()[idx+length:])
						}
						if foo.AdvanceLine {
							p.nextLine()
						}
					}
					return token, nil
				}
			}
		}
		c, err := p.nextRune()
		fmt.Printf("Line %d, offset %d: %#U\n", p.Position.Line, p.Position.LineOffset, c)
		if err != nil {
			//if err != io.EOF {
			return nil, err
			//}
			break
		}
		p.buf.WriteRune(c)
	}
	return nil, nil
}
