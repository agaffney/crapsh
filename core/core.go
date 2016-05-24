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
	} else if state.config.FileProvided {
		file, err := os.Open(state.config.File)
		if err != nil {
			fmt.Printf("%s: %s\n", state.config.Binary, err)
		}
		state.config.Binary = state.config.File
		state.parser.Parse(file)
	} else {
		// Code to show prompt
	}
	for line := range state.parser.LineChan {
		util.DumpJson(line)
	}
	if state.parser.Error != nil {
		fmt.Printf("%s: %s\n", state.config.Binary, state.parser.Error)
	}
}
