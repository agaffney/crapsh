package builtin

import (
	"fmt"
	"github.com/agaffney/crapsh/core/state"
)

func Help(state *state.State, args []string) int {
	if len(args) == 1 {
		fmt.Printf("Usage: help <topic>\n")
		return 0
	} else {
		for _, builtin := range Builtins {
			if builtin.Name == args[1] {
				if builtin.HelpText == `` {
					fmt.Printf("help: %s: no help text available\n", args[1])
				} else {
					fmt.Printf(builtin.HelpText)
				}
				return 0
			}
		}
	}
	fmt.Printf("help: %s: no such topic\n", args[1])
	return 1
}

func init() {
	registerBuiltin(Builtin{Name: `help`, Entrypoint: Help})
}
