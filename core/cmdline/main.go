package cmdline

import (
	"github.com/agaffney/crapsh/core"
	"github.com/agaffney/crapsh/core/cmdline/parser"
	"github.com/agaffney/crapsh/core/flags"
	"os"
	"path"
)

func Main() {
	c := &core.Config{}
	c.Binary = path.Base(os.Args[0])
	parseCmdlineOpts(c)
	state := core.NewState(c)
	state.Run()
}

func parseCmdlineOpts(c *core.Config) error {
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
	inputOptions.Add([]*parser.Option{
		{Long: `help`},
	})
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
