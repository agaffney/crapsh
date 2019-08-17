package lexer

import (
	"bytes"
	parser_input "github.com/agaffney/crapsh/parser/input"
	"io"
)

const BUF_THRESHOLD = 1024

type Lexer struct {
	buf            *bytes.Buffer
	tokenChan      chan *Token
	errorChan      chan error
	input          parser_input.Input
	lineNum        int
	lineOffset     int
	prevLineOffset int
}

type Token struct {
	Type    int
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

func (l *Lexer) Start(input parser_input.Input) {
	l.Reset()
	l.input = input
	// TODO: check for error
	l.readLine(false)
	go l.generateTokens()
}

func (l *Lexer) readLine(continuation bool) error {
	line, err := l.input.ReadLine(continuation)
	if line != "" {
		l.buf.WriteString(line)
	}
	return err
}

func (l *Lexer) GetError() error {
	e, ok := <-l.errorChan
	if ok {
		return e
	} else {
		return nil
	}
}

func (l *Lexer) ReadToken() (*Token, error) {
	t, ok := <-l.tokenChan
	if ok {
		return t, nil
	} else {
		return nil, io.EOF
	}
}

// Return a single character (rune) from the buffer
func (l *Lexer) nextRune() (rune, error) {
	r, _, err := l.buf.ReadRune()
	// Preserve previous line offset in case we need to unread a rune
	l.prevLineOffset = l.lineOffset
	l.lineOffset++
	return r, err
}

// Rewind buffer position by one character (rune)
func (l *Lexer) unreadRune() error {
	err := l.buf.UnreadRune()
	if l.lineOffset == 0 {
		l.lineNum--
	}
	l.lineOffset = l.prevLineOffset
	return err
}

// Increment line number for input
func (l *Lexer) nextLine() {
	l.lineNum++
	l.lineOffset = 0
}
