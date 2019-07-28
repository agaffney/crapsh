package lexer

import (
	"bufio"
	"bytes"
	"io"
)

const BUF_THRESHOLD = 1024

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

func (l *Lexer) GetError() error {
	e, ok := <-l.errorChan
	if ok {
		return e
	} else {
		return nil
	}
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
	for {
		// Check the buffer at the beginning to catch tokens already in the buffer
		// from the last iteration
		if l.buf.Len() > 0 {
			for _, foo := range TokenDefinitions {
				if ok, value := foo.Match(l.buf, 0); ok {
					token = &Token{Type: foo.Name, Value: value, LineNum: l.lineNum, Offset: l.lineOffset}
					l.lineOffset += len(value)
					if len(value) == len(l.buf.String()) {
						l.buf.Reset()
					} else {
						l.buf = bytes.NewBufferString(l.buf.String()[len(value):])
					}
					if foo.AdvanceLine {
						l.nextLine()
					}
					l.tokenChan <- token
					break
				}
			}
		}
		if l.buf.Len() < BUF_THRESHOLD {
			for i := 0; i < BUF_THRESHOLD; i++ {
				r, _, err := l.input.ReadRune()
				if err != nil {
					if err != io.EOF {
						l.errorChan <- err
						break
					}
				} else {
					l.buf.WriteRune(r)
				}
			}
			if l.buf.Len() == 0 {
				break
			}
		}
	}
	close(l.tokenChan)
}
