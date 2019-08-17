package input

import (
	"github.com/chzyer/readline"
	"io"
	"os"
)

type Interactive struct {
	readline *readline.Instance
}

func NewInteractive() *Interactive {
	i := &Interactive{}
	rl, err := readline.NewEx(&readline.Config{
		HistoryFile: "/tmp/readline.tmp",
		//AutoComplete:    completer,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",

		HistorySearchFold: true,
		//FuncFilterInputRune: filterInput,
	})
	i.readline = rl
	if err != nil {
		panic(err)
	}
	//defer rl.Close()
	return i
}

func (i *Interactive) IsAvailable() bool {
	return true
}

func (i *Interactive) ReadLine(continuation bool) (string, error) {
	if continuation {
		// TODO: pull from $PS2
		i.readline.SetPrompt("> ")
	} else {
		// TODO: pull from $PS1
		i.readline.SetPrompt("\033[31m»\033[0m ")
	}
	line, err := i.readline.Readline()
	// TODO: do this elsewhere
	// Hack to allow us to exit cleanly
	if line == "logout" || line == "exit" {
		os.Exit(0)
	}
	if err != nil {
		if err == readline.ErrInterrupt {
			return "", nil
		} else if err == io.EOF {
			os.Exit(0)
		}
		return "", err
	}
	// Append a newline for consistency with other input methods
	line += "\n"
	return line, nil
}
