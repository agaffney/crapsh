package core

import (
	"fmt"
	"github.com/agaffney/crapsh/parser"
	"github.com/agaffney/crapsh/util"
	"os"
	"strings"
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
		state.parser.Parse(strings.NewReader(state.config.Command))
	} else if state.config.FileProvided && !state.config.ReadFromStdin {
		file, err := os.Open(state.config.File)
		if err != nil {
			fmt.Printf("%s: %s\n", state.config.Binary, err)
		}
		state.config.Binary = state.config.File
		state.parser.Parse(file)
	} else {
		// Code to show prompt
		fmt.Println("Interactive prompt not currently supported")
		os.Exit(1)
	}
	for {
		cmd := state.parser.GetCommand()
		if cmd == nil {
			fmt.Println("no more commands")
			break
		}
		util.DumpJson(cmd, "Command:\n")
	}
	if err := state.parser.GetError(); err != nil {
		fmt.Printf("%s: %s\n", state.config.Binary, err)
		os.Exit(1)
	}
	os.Exit(0)
}
