package input

import (
	"bufio"
	"strings"
)

type StringParserInput struct {
	input *bufio.Reader
}

func NewStringParserInput(input string) *StringParserInput {
	i := &StringParserInput{}
	i.input = bufio.NewReader(strings.NewReader(input))
	return i
}

func (i *StringParserInput) ReadLine(continuation bool) (string, error) {
	return i.input.ReadString('\n')
}

func (i *StringParserInput) IsAvailable() bool {
	return (i.input.Buffered() > 0)
}
