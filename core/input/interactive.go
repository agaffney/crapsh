package input

import (
	// TODO: maybe replace with https://github.com/fiorix/go-readline
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
		EOFPrompt:       "logout",

		HistorySearchFold: true,
		//FuncFilterInputRune: filterInput,
	})
	i.readline = rl
	if err != nil {
		panic(err)
	}
	// TODO: figure out where to close the readline instance on exit
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
		i.readline.SetPrompt("$ ") // \033[31m»\033[0m ")
	}
	line, err := i.readline.Readline()
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
