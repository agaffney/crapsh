package lexer

import (
	"bufio"
	"bytes"
	"io"
)

type Lexer struct {
	buf       *bytes.Buffer
	tokenChan chan *Token
	input     io.Reader
	lineNum int
	lineOffset int
}

type Token struct {
	LineNum int
	Offset  int
	Value   string
}

func New() *Lexer {
	l := &Lexer{}
	l.buf = bytes.NewBuffer(nil)
	l.tokenChan = make(chan *Token, 100)
	return l
}

func (l *Lexer) Reset() {

}

func (l *Lexer) Start(input io.Reader) {
	l.input := bufio.NewReader(input)
	go func() {
		for {
			r, _, err := l.input.ReadRune()
			if err == io.EOF {
				break
			}
			if err != nil {
				// TODO: handle other errors
			}
			// TODO: use token definitions
			tok := &Token{Value: string(r)}
			l.tokenChan <- tok
		}
		close(l.tokenChan)
	}()
}

func (l *Lexer) ReadToken() *Token {
	t, ok := <-l.tokenChan
	if ok {
		return t
	} else {
		return nil
	}
}

// Return a single character (rune) from the buffer
func (l *Lexer) nextRune() (rune, error) {
	r, _, err := l.input.ReadRune()
	l.LineOffset++
	return r, err
}

// Rewind buffer position by one character (rune)
func (l *Lexer) unreadRune() error {
	err := l.input.UnreadRune()
	l.LineOffset--
	return err
}

// Increment line number for input
func (l *Lexer) nextLine() {
	l.lineNum++
	l.LineOffset = 0
}

/*
// Scan input buffer for a matching token
func (l *Lexer) generateTokens(input io.Reader) (error) {
	var buf_len int

	for {
		// Check the buffer at the beginning to catch tokens already in the buffer
		// from the last iteration
		buf_len = p.buf.Len()
		if buf_len > 0 {
			for _, foo := range tokens.Tokens {
				if idx, length := foo.Match(p.buf); idx > -1 {
					//fmt.Printf("idx=%d, length=%d, data='%s'\n", idx, length, p.buf.Bytes()[idx:idx+length])
					if length == buf_len && foo.MatchUntilNextToken {
						break
					}
					var token *Token
					if idx > 0 {
						// Return data up to token as "generic" token and remove from buffer
						token = &Token{Type: `Generic`, Value: string(p.buf.Bytes()[0:idx]), Pos: p.Position}
						p.buf = bytes.NewBuffer(p.buf.Bytes()[idx:])
					} else {
						token = &Token{Type: foo.Name, Value: string(p.buf.Bytes()[idx : idx+length]), Pos: p.Position}
						if length == p.buf.Len() {
							p.buf.Reset()
						} else {
							p.buf = bytes.NewBuffer(p.buf.Bytes()[idx+length:])
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
*/
