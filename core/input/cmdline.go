package input

import (
	"bufio"
	"strings"
)

type Cmdline struct {
	input *bufio.Reader
}

func NewCmdline(input string) *Cmdline {
	i := &Cmdline{input: bufio.NewReader(strings.NewReader(input))}
	return i
}

func (i *Cmdline) ReadLine(continuation bool) (string, error) {
	return i.input.ReadString('\n')
}

func (i *Cmdline) IsAvailable() bool {
	return (i.input.Buffered() > 0)
}
