package core

import (
	"fmt"
	"github.com/agaffney/crapsh/cmd/builtin"
	"github.com/agaffney/crapsh/core/config"
	//"github.com/agaffney/crapsh/core/flags"
	core_input "github.com/agaffney/crapsh/core/input"
	"github.com/agaffney/crapsh/core/state"
	"github.com/agaffney/crapsh/parser"
	parser_input "github.com/agaffney/crapsh/parser/input"
	//"github.com/agaffney/crapsh/util"
	"os"
)

type Core struct {
	config *config.Config
	state  *state.State
}

func New(config *config.Config) *Core {
	core := &Core{config: config}
	core.state = state.New()
	return core
}

func (core *Core) Run() {
	var input parser_input.Input
	if core.config.CommandProvided {
		// Command provided via -c option
		input = core_input.NewCmdline(core.config.Command)
	} else if core.config.FileProvided && !core.config.ReadFromStdin {
		// Read commands from STDIN (-s option)
		/*
			// TODO: move to core/input/file.go
			file, err := os.Open(core.config.File)
			if err != nil {
				fmt.Printf("%s: %s\n", core.config.Binary, err)
			}
			core.config.Binary = core.config.File
			core.state.parser.Parse(file)
			core.state.processCommands()
		*/
	} else {
		// Interactive input
		input = core_input.NewInteractive()
	}
	core.state.Parser = parser.NewParser(input)
	core.processCommands()
	os.Exit(0)
}

func (core *Core) processCommands() {
	for {
		cmd, err := core.state.Parser.GetCommand()
		if err != nil {
			fmt.Printf("%s: %s\n", core.config.Binary, err.Error())
			os.Exit(1)
		}
		if cmd == nil {
			fmt.Println("no more commands")
			break
		}
		//util.DumpJson(cmd, "Command:\n")
		for _, pipeline := range cmd.GetChildren() {
			for _, command := range pipeline.GetChildren() {
				args := []string{}
				for _, node := range command.GetChildren() {
					if node.GetName() == `Word` {
						args = append(args, node.GetToken().Value)
					}
				}
				fmt.Printf("Args: %#v\n", args)
				foundBuiltin := false
				for _, b := range builtin.Builtins {
					if b.Name == args[0] {
						foundBuiltin = true
						ret := b.Entrypoint(core.state, args)
						fmt.Printf("Returned %d\n", ret)
						break
					}
				}
				if !foundBuiltin {
					fmt.Printf("Command not found\n")
				}
			}
		}
	}
}
