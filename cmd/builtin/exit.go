package builtin

import (
	"github.com/agaffney/crapsh/core/state"
	"os"
)

func Exit(state *state.State, args []string) int {
	os.Exit(0)
	return 0
}

func init() {
	registerBuiltin(Builtin{Name: `exit`, Entrypoint: Exit})
	registerBuiltin(Builtin{Name: `logout`, Entrypoint: Exit})
}
