package builtin

import (
	"fmt"
	"github.com/agaffney/crapsh/core/state"
)

func Help(state *state.State, args []string) int {
	fmt.Printf("I can only help those that help themselves\n")
	return 0
}

func init() {
	registerBuiltin(Builtin{Name: `help`, Entrypoint: Help})
}
