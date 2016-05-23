package cmd

import (
	"github.com/agaffney/crapsh/core"
	//"github.com/agaffney/crapsh/prompt"
)

func Start() {
	//fmt.Println("in main()")
	//p := &prompt.Prompt{
	//	Text: "prompt text",
	//}
	//p.Show()
	c := &core.Config{}
	parse_cmdline_opts(c)
	state := core.New(c)
	state.Start()
}
