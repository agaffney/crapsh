package cmd

import (
	"fmt"
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
	fmt.Printf("Command is '%s'\n", c.Command)
}
