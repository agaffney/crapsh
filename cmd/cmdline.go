package cmd

import (
	"flag"
	"github.com/agaffney/crapsh/core"
)

type StringBool struct {
	value string
	used  bool
}

func (s *StringBool) Set(value string) error {
	s.value = value
	s.used = true
	return nil
}

func (s *StringBool) String() string {
	return s.value
}

func parse_cmdline_opts(c *core.Config) {
	flag_command := &StringBool{}
	flag.Var(flag_command, "c", "specifies a command to run non-interactively")
	flag.BoolVar(&c.ReadFromStdin, "s", false, "read commands from STDIN (default if no file provided)")
	flag.Parse()
	if flag_command.used {
		c.Command = flag_command.String()
		c.CommandProvided = true
	}
	c.Args = flag.Args()
	if len(c.Args) > 0 {
		c.File = c.Args[0]
		c.FileProvided = true
	}
}
