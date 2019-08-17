package input

import (
	"bufio"
	"strings"
)

type Input interface {
	ReadLine() (string, error)
	ReadAnotherLine() (string, error)
}

type StringParserInput struct {
	input *bufio.Reader
}

func NewStringParserInput(input string) *StringParserInput {
	i := &StringParserInput{}
	i.input = bufio.NewReader(strings.NewReader(input))
	return i
}

func (i *StringParserInput) ReadLine() (string, error) {
	return i.input.ReadString('\n')
}

func (i *StringParserInput) ReadAnotherLine() (string, error) {
	return i.input.ReadString('\n')
}
