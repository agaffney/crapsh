package parser

import (
	"fmt"
	"os"
)

const (
	_               = iota
	TYPE_FLAG       // simple on/off
	TYPE_SHELL_FLAG // flag corresponding to 'set' option
	TYPE_ARG        // flag that expects an argument
)

type Option struct {
	Type  int
	Short string
	Long  string
	Set   bool // whether the flag was specified
	Value bool
	Arg   string
}

type OptionSet struct {
	options []*Option
}

func (o *OptionSet) Add(options []*Option) {
	o.options = append(o.options, options...)
}

func (o *OptionSet) Options() []*Option {
	return o.options
}

func (o OptionSet) FindOption(flag string, shellFlag bool) *Option {
	flagLen := len(flag)
	for _, option := range o.options {
		if flagLen == 1 && option.Short == flag {
			if !shellFlag || option.Type == TYPE_SHELL_FLAG {
				return option
			}
		} else if flagLen > 1 && option.Long == flag {
			if !shellFlag || option.Type == TYPE_SHELL_FLAG {
				return option
			}
		}
	}
	return nil
}

func copyOption(src *Option) *Option {
	newOption := &Option{}
	newOption.Type = src.Type
	newOption.Short = src.Short
	newOption.Long = src.Long
	newOption.Set = src.Set
	newOption.Value = src.Value
	newOption.Arg = src.Arg
	return newOption
}

func Parse(options OptionSet) (*OptionSet, []string, error) {
	newOptionSet := &OptionSet{}
	doneWithOptions := false
	expectingArg := false
	var prevOption *Option
	var prevOptionDisplay string
	args := []string{}
	for _, arg := range os.Args[1:] {
		if doneWithOptions {
			args = append(args, arg)
			continue
		}
		if arg[0:2] == `--` {
			if len(arg) == 2 {
				doneWithOptions = true
				continue
			}
			option := options.FindOption(arg[2:], false)
			if option == nil {
				return nil, nil, fmt.Errorf("unknown option: %s", arg)
			}
			newOption := copyOption(option)
			newOption.Set = true
			if newOption.Type == TYPE_ARG {
				expectingArg = true
				prevOption = newOption
				prevOptionDisplay = arg
			}
			newOptionSet.Add([]*Option{newOption})
		} else if arg[0] == '-' || arg[0] == '+' {
			if len(arg) == 1 {
				args = append(args, arg)
				expectingArg = false
				continue
			}
			if expectingArg {
				return nil, nil, fmt.Errorf("%s: option requires an argument", prevOptionDisplay)
			}
			shellFlag := false
			if arg[0] == '+' {
				shellFlag = true
			}
			for _, flag := range arg[1:] {
				option := options.FindOption(string(flag), shellFlag)
				if option == nil {
					// Append to args if it's a lone + flag that doesn't match an option
					if arg[0] == '+' && len(arg) == 2 {
						args = append(args, arg)
						continue
					}
					return nil, nil, fmt.Errorf("unknown option: %c%c", arg[0], flag)
				}
				foundExisting := false
				newOption := newOptionSet.FindOption(string(flag), false)
				if newOption == nil {
					newOption = copyOption(option)
				} else {
					foundExisting = true
				}
				newOption.Set = true
				if arg[0] == '-' {
					newOption.Value = true
				} else {
					newOption.Value = false
				}
				if newOption.Type == TYPE_ARG {
					expectingArg = true
					prevOption = newOption
					prevOptionDisplay = fmt.Sprintf("%c%c", arg[0], flag)
				}
				if !foundExisting {
					newOptionSet.Add([]*Option{newOption})
				}
			}
		} else {
			if expectingArg {
				expectingArg = false
				prevOption.Arg = arg
			} else {
				args = append(args, arg)
			}
		}
	}
	if expectingArg {
		return nil, nil, fmt.Errorf("%s: option requires an argument", prevOptionDisplay)
	}
	return newOptionSet, args, nil
}
