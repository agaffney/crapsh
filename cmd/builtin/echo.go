package builtin

import (
	"fmt"
	"github.com/agaffney/crapsh/core/cmdline/parser"
	"github.com/agaffney/crapsh/core/state"
)

func Echo(state *state.State, inputArgs []string) int {
	inputOptions := parser.OptionSet{}
	inputOptions.Add([]*parser.Option{
		{Short: `n`},
		{Short: `e`},
		{Short: `E`},
	})
	options, args, err := parser.Parse(inputOptions, inputArgs[1:])
	if err != nil {
		fmt.Printf("echo: %s\n", err.Error())
		return 1
	}
	appendNewline := true
	enableEscapes := false
	if option := options.FindOption(`n`, false); option.Set {
		appendNewline = false
	}
	if option := options.FindOption(`e`, false); option.Set {
		enableEscapes = true
	}
	if option := options.FindOption(`E`, false); option.Set {
		enableEscapes = false
	}
	for idx, arg := range args {
		if idx > 0 {
			fmt.Print(" ")
		}
		if enableEscapes {
			// TODO: do something
		}
		fmt.Print(arg)
	}
	if appendNewline {
		fmt.Print("\n")
	} else {
		// Until https://github.com/chzyer/readline/issues/169 is fixed
		fmt.Print("\n")
	}
	return 0
}

var EchoHelpText = `Usage: echo [-neE] [arg ...]

    Write arguments to the standard output.

    Display the ARGs, separated by a single space character and followed by a
    newline, on the standard output.

    Options:
      -n        do not append a newline
      -e        enable interpretation of the following backslash escapes
      -E        explicitly suppress interpretation of backslash escapes

    'echo' interprets the following backslash-escaped characters:
      \a        alert (bell)
      \b        backspace
      \c        suppress further output
      \e        escape character
      \E        escape character
      \f        form feed
      \n        new line
      \r        carriage return
      \t        horizontal tab
      \v        vertical tab
      \\        backslash
      \0nnn     the character whose ASCII code is NNN (octal).  NNN can be
        0 to 3 octal digits
      \xHH      the eight-bit character whose value is HH (hexadecimal).  HH
        can be one or two hex digits

    Exit Status:
    Returns success unless a write error occurs.
`

func init() {
	registerBuiltin(Builtin{Name: `echo`, Entrypoint: Echo, HelpText: EchoHelpText})
}
