package cmd

import (
	"github.com/agaffney/crapsh/core"
	//"github.com/agaffney/crapsh/prompt"
	"os"
)

func Start() {
	//fmt.Println("in main()")
	//p := &prompt.Prompt{
	//	Text: "prompt text",
	//}
	//p.Show()
	c := &core.Config{}
	c.Binary = os.Args[0]
	parse_cmdline_opts(c)
	state := core.New(c)
	state.Start()
}
