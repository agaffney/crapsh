package core

import (
	"github.com/agaffney/crapsh/parser"
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
	state.parser.Parse(state.config.Command)
}
