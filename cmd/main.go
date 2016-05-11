package cmd

import (
	"fmt"
	"github.com/agaffney/crapsh/parser"
	"github.com/agaffney/crapsh/prompt"
)

func Start() {
	fmt.Println("in main()")
	p := &prompt.Prompt{
		Text: "prompt text",
	}
	p.Show()
	c := &Config{}
	parse_cmdline_opts(c)
	parser := parser.NewParser()
	parser.Parse(c.Command)
}
