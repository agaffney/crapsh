package cmd

import (
	"flag"
)

func parse_cmdline_opts(c *Config) {
	flag.StringVar(&c.Command, "c", "", "specifies a command to run non-interactively")
	flag.Parse()
	c.Args = flag.Args()
}
