package cmd

import (
	"flag"
	"github.com/agaffney/crapsh/core"
)

func parse_cmdline_opts(c *core.Config) {
	flag.StringVar(&c.Command, "c", "", "specifies a command to run non-interactively")
	flag.Parse()
	c.Args = flag.Args()
}
