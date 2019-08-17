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
		input := core_input.NewCmdline(state.config.Command)
		state.parser.Parse(input)
		state.processCommands()
		/*
			} else if state.config.FileProvided && !state.config.ReadFromStdin {
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
		input := core_input.NewInteractive()
		state.parser.Parse(input)
		state.processCommands()
	}
	os.Exit(0)
}

func (state *State) processCommands() {
	for {
		cmd := state.parser.GetCommand()
		if cmd == nil {
			fmt.Println("no more commands")
			break
		}
		util.DumpJson(cmd, "Command:\n")
	}
	/*
		if err := state.parser.GetError(); err != nil {
			fmt.Printf("%s: %s\n", state.config.Binary, err)
			os.Exit(1)
		}
	*/
}
