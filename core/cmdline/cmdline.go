package cmdline

import (
	"github.com/agaffney/crapsh/core"
	"github.com/agaffney/crapsh/core/cmdline/parser"
	"github.com/agaffney/crapsh/core/flags"
	//"github.com/agaffney/crapsh/util"
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

func parse_cmdline_opts(c *core.Config) error {
	inputOptions := parser.OptionSet{}
	for _, flag := range flags.Flags {
		if flag.Short == `` && !flag.CmdlineOnly {
			continue
		}
		option := &parser.Option{Short: flag.Short, Long: flag.Long, Type: parser.TYPE_SHELL_FLAG}
		if !flag.CmdlineOnly {
			option.Type = parser.TYPE_FLAG
		} else if flag.HasArg {
			option.Type = parser.TYPE_ARG
		}
		inputOptions.Add([]*parser.Option{option})
	}
	options, args, err := parser.Parse(inputOptions)
	if err != nil {
		return err
	}
	optCommand := options.FindOption("c", false)
	if optCommand.Set {
		c.Command = optCommand.Arg
		c.CommandProvided = true
	}
	c.Args = args
	if len(c.Args) > 0 {
		c.File = c.Args[0]
		c.FileProvided = true
	}
	return nil
}
