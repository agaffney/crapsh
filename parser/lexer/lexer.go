package lexer

import (
	"bufio"
	"bytes"
	"io"
)

type Lexer struct {
	buf       *bytes.Buffer
	tokenChan chan *Token
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
	go func() {
		reader := bufio.NewReader(input)
		for {
			r, _, err := reader.ReadRune()
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
