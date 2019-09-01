package builtin

import (
	"github.com/agaffney/crapsh/core/state"
)

type Builtin struct {
	Name       string
	Entrypoint EntrypointFunc
}

type EntrypointFunc func(*state.State, []string) int

var Builtins []Builtin

func registerBuiltin(b Builtin) {
	Builtins = append(Builtins, b)
}
