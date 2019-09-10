package core

import (
	"fmt"
	"github.com/agaffney/crapsh/core/config"
	"github.com/agaffney/crapsh/core/executor"
	core_input "github.com/agaffney/crapsh/core/input"
	"github.com/agaffney/crapsh/core/state"
	"github.com/agaffney/crapsh/parser"
	parser_input "github.com/agaffney/crapsh/parser/input"
	//"github.com/agaffney/crapsh/util"
	"os"
)

type Core struct {
	config   *config.Config
	state    *state.State
	executor *executor.Executor
}

func New(config *config.Config) *Core {
	core := &Core{config: config}
	core.state = state.New()
	core.executor = executor.New()
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
		cmdAst, err := core.state.Parser.GetCommand()
		if err != nil {
			fmt.Printf("%s: %s\n", core.config.Binary, err.Error())
			os.Exit(1)
		}
		if cmdAst == nil {
			//fmt.Println("no more commands")
			break
		}
		//util.DumpJson(cmdAst, "Command (AST):\n")
		completeCommand := executor.NewCompleteCommand(cmdAst)
		//util.DumpJson(completeCommand, "Complete command (executor):\n")
		err = completeCommand.Run(core.state)
		if err != nil {
			fmt.Printf("%s: %s\n", core.config.Binary, err)
		}
	}
}
