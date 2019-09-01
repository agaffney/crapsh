package cmdline

import (
	"github.com/agaffney/crapsh/core"
	"os"
	"path"
)

func Main() {
	c := &core.Config{}
	c.Binary = path.Base(os.Args[0])
	parse_cmdline_opts(c)
	state := core.New(c)
	state.Run()
}
