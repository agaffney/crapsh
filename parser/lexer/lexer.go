package lexer

import (
	"bytes"
	parser_input "github.com/agaffney/crapsh/parser/input"
)

const (
	BUF_THRESHOLD = 1024
	MODE_NULL     = iota
	MODE_NORMAL
	MODE_HEREDOC
)

type Lexer struct {
	lineBuf        *bytes.Buffer
	input          parser_input.Input
	lineNum        int
	lineOffset     int
	prevLineOffset int
	mode           int
}

type Token struct {
	Type    int
	LineNum int
	Offset  int
	Value   string
}

func New(input parser_input.Input) *Lexer {
	l := &Lexer{input: input}
	l.Reset()
	return l
}

func (l *Lexer) Reset() {
	l.lineBuf = bytes.NewBuffer(nil)
	l.lineNum = 1
	l.lineOffset = 1
	l.mode = MODE_NORMAL
}

func (l *Lexer) SetMode(mode int) {
	l.mode = mode
}

func (l *Lexer) ReadLine() error {
	return l.readLine(false)
}

func (l *Lexer) readLine(continuation bool) error {
	line, err := l.input.ReadLine(continuation)
	if line != "" {
		l.lineBuf.WriteString(line)
	}
	return err
}

func (l *Lexer) ReadToken() (*Token, error) {
	return l.nextToken()
}

// Return a single character (rune) from the buffer
func (l *Lexer) nextRune() (rune, error) {
	r, _, err := l.lineBuf.ReadRune()
	// Preserve previous line offset in case we need to unread a rune
	l.prevLineOffset = l.lineOffset
	l.lineOffset++
	return r, err
}

// Rewind buffer position by one character (rune)
func (l *Lexer) unreadRune() error {
	err := l.lineBuf.UnreadRune()
	if l.lineOffset == 1 {
		l.lineNum--
	}
	l.lineOffset = l.prevLineOffset
	return err
}

// Increment line number for input
func (l *Lexer) nextLine() {
	l.lineNum++
	l.lineOffset = 1
}
