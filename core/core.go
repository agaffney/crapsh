package core

import (
	"fmt"
	core_input "github.com/agaffney/crapsh/core/input"
	"github.com/agaffney/crapsh/parser"
	"github.com/agaffney/crapsh/util"
	"os"
)

type State struct {
	parser *parser.Parser
	config *Config
}

func New(config *Config) *State {
	state := &State{config: config, parser: parser.NewParser()}
	return state
}

func (state *State) Start() {
	if state.config.CommandProvided {
		// Command provided via -c option
		input := core_input.NewCmdline(state.config.Command)
		state.parser.Start(input)
	} else if state.config.FileProvided && !state.config.ReadFromStdin {
		// Read commands from STDIN (-s option)
		/*
			// TODO: move to core/input/file.go
			file, err := os.Open(state.config.File)
			if err != nil {
				fmt.Printf("%s: %s\n", state.config.Binary, err)
			}
			state.config.Binary = state.config.File
			state.parser.Parse(file)
			state.processCommands()
		*/
	} else {
		// Interactive input
		input := core_input.NewInteractive()
		state.parser.Start(input)
	}
	state.processCommands()
	os.Exit(0)
}

func (state *State) processCommands() {
	for {
		cmd, err := state.parser.GetCommand()
		if err != nil {
			fmt.Printf("%s: %s\n", state.config.Binary, err.Error())
			os.Exit(1)
		}
		if cmd == nil {
			fmt.Println("no more commands")
			break
		}
		util.DumpJson(cmd, "Command:\n")
	}
}
