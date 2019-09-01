package state

import (
	"github.com/agaffney/crapsh/core/flags"
	"github.com/agaffney/crapsh/parser"
)

type State struct {
	Parser *parser.Parser
	Flags  []FlagState
}

type FlagState struct {
	flags.Flag
	Set bool
}

func New() *State {
	state := &State{}
	return state
}
