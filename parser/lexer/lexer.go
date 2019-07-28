package lexer

import (
	"bufio"
	"bytes"
	//"fmt"
	"io"
)

type Lexer struct {
	buf        *bytes.Buffer
	tokenChan  chan *Token
	errorChan  chan error
	input      *bufio.Reader
	lineNum    int
	lineOffset int
}

type Token struct {
	Type    string
	LineNum int
	Offset  int
	Value   string
}

func New() *Lexer {
	l := &Lexer{}
	l.Reset()
	return l
}

func (l *Lexer) Reset() {
	l.buf = bytes.NewBuffer(nil)
	l.tokenChan = make(chan *Token, 100)
	l.errorChan = make(chan error, 1)
	l.input = nil
	l.lineNum = 1
	l.lineOffset = 1
}

func (l *Lexer) Start(input io.Reader) {
	l.input = bufio.NewReader(input)
	go l.generateTokens()
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
	l.lineOffset++
	return r, err
}

// Rewind buffer position by one character (rune)
func (l *Lexer) unreadRune() error {
	err := l.input.UnreadRune()
	l.lineOffset--
	return err
}

// Increment line number for input
func (l *Lexer) nextLine() {
	l.lineNum++
	l.lineOffset = 0
}

// Scan input buffer for a matching token
func (l *Lexer) generateTokens() {
	var token *Token
	var foundEOF bool
	for {
		// Check the buffer at the beginning to catch tokens already in the buffer
		// from the last iteration
		if l.buf.Len() > 0 {
			for _, foo := range TokenDefinitions {
				if idx, length := foo.Match(l.buf); idx > -1 {
					//fmt.Printf("idx=%d, length=%d, data='%s'\n", idx, length, l.buf.Bytes()[idx:idx+length])
					if length == l.buf.Len() && foo.MatchUntilNextToken {
						break
					}
					if idx > 0 {
						// Return data up to token as "generic" token and remove from buffer
						token = &Token{Type: `Generic`, Value: string(l.buf.Bytes()[0:idx]), LineNum: l.lineNum, Offset: l.lineOffset}
						l.lineOffset += idx
						l.buf = bytes.NewBuffer(l.buf.Bytes()[idx:])
					} else {
						token = &Token{Type: foo.Name, Value: string(l.buf.Bytes()[idx : idx+length]), LineNum: l.lineNum, Offset: l.lineOffset}
						l.lineOffset += length
						if length == l.buf.Len() {
							l.buf.Reset()
						} else {
							l.buf = bytes.NewBuffer(l.buf.Bytes()[idx+length:])
						}
						if foo.AdvanceLine {
							l.nextLine()
						}
					}
					l.tokenChan <- token
					break
				}
			}
			if foundEOF && l.buf.Len() > 0 {
				// Return remaining data as "generic" token
				token = &Token{Type: `Generic`, Value: l.buf.String(), LineNum: l.lineNum, Offset: l.lineOffset}
				l.buf.Reset()
				l.tokenChan <- token
			}
		}
		r, _, err := l.input.ReadRune()
		if err != nil {
			if err != io.EOF {
				l.errorChan <- err
				break
			} else {
				foundEOF = true
			}
		} else {
			l.buf.WriteRune(r)
		}
		if l.buf.Len() == 0 {
			break
		}
	}
	close(l.tokenChan)
}
