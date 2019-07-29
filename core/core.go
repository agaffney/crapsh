package core

import (
	"bytes"
	"fmt"
	"github.com/agaffney/crapsh/parser"
	"github.com/agaffney/crapsh/util"
	"github.com/chzyer/readline"
	"io"
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
		state.processCommands()
	} else if state.config.FileProvided && !state.config.ReadFromStdin {
		file, err := os.Open(state.config.File)
		if err != nil {
			fmt.Printf("%s: %s\n", state.config.Binary, err)
		}
		state.config.Binary = state.config.File
		state.parser.Parse(file)
		state.processCommands()
	} else {
		rl, err := readline.NewEx(&readline.Config{
			Prompt:      "\033[31m»\033[0m ",
			HistoryFile: "/tmp/readline.tmp",
			//AutoComplete:    completer,
			InterruptPrompt: "^C",
			EOFPrompt:       "exit",

			HistorySearchFold: true,
			//FuncFilterInputRune: filterInput,
		})
		if err != nil {
			panic(err)
		}
		defer rl.Close()
		for {
			line, err := rl.Readline()
			if err == readline.ErrInterrupt {
				if len(line) == 0 {
					break
				} else {
					continue
				}
			} else if err == io.EOF {
				break
			}
			fmt.Println(line)
			buf := bytes.NewBufferString(line)
			state.parser.Parse(buf)
			state.processCommands()
		}
	}
	os.Exit(0)
}

func (state *State) processCommands() {
	for {
		cmd := state.parser.GetCommand()
		if cmd == nil {
			fmt.Println("no more commands")
			break
		}
		util.DumpJson(cmd, "Command:\n")
	}
	if err := state.parser.GetError(); err != nil {
		fmt.Printf("%s: %s\n", state.config.Binary, err)
		os.Exit(1)
	}
}
