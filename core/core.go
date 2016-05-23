package core

import (
	"github.com/agaffney/crapsh/parser"
	"github.com/agaffney/crapsh/util"
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
	state.parser.Parse(strings.NewReader(state.config.Command))
	for line := range state.parser.LineChan {
		util.DumpJson(line)
	}
}
