package parser

import (
	"os"
	"testing"
)

type parserTestCaseFlag struct {
	flag  string
	value bool
	arg   string
}

type parserTestCase struct {
	input   []string
	options OptionSet
	flags   []parserTestCaseFlag
	args    []string
}

func runTests(testCases []parserTestCase, t *testing.T) {
	for _, testCase := range testCases {
		os.Args = []string{"test"}
		os.Args = append(os.Args, testCase.input...)
		options, args, err := Parse(testCase.options)
		if err != nil {
			t.Fatalf("unexpected error: %s", err.Error())
		}
		argsMatch := true
		if len(args) != len(testCase.args) {
			argsMatch = false
		}
		for idx, arg := range testCase.args {
			if arg != args[idx] {
				argsMatch = false
				break
			}
		}
		if !argsMatch {
			t.Fatalf("expected args: %v, got: %v", testCase.args, args)
		}
		// Check that all expected flags are set
		for _, expected := range testCase.flags {
			option := options.FindOption(expected.flag, false)
			if !option.Set {
				t.Fatalf("option '%s' did not have expected Set value, got: %t", expected.flag, option.Set)
			}
			if option.Type == TYPE_ARG {
				if option.Arg != expected.arg {
					t.Fatalf("option '%s' did not have expected Arg. expected: %s, got %s", expected.flag, expected.arg, option.Arg)
				}
			} else {
				if option.Value != expected.value {
					t.Fatalf("option '%s' did not have expected Value, got: %t", expected.flag, option.Value)
				}
			}
		}
		// Check for unexpected flags to be set
		for _, option := range options.Options() {
			if option.Set {
				foundFlag := false
				for _, flag := range testCase.flags {
					if flag.flag == option.Short || flag.flag == option.Long {
						foundFlag = true
						break
					}
				}
				if !foundFlag {
					t.Fatalf("option '%s' was unexpectedly set", option.Short)
				}
			}
		}
	}
}

func TestParserBasic(t *testing.T) {
	commonOptions := OptionSet{}
	commonOptions.Add([]*Option{
		{Short: `c`, Type: TYPE_ARG},
		{Short: `x`},
	})
	testCases := []parserTestCase{
		{
			input:   []string{`-c`, `foo bar baz`},
			options: commonOptions,
			flags: []parserTestCaseFlag{
				{
					flag: `c`,
					arg:  `foo bar baz`,
				},
			},
		},
		{
			input:   []string{`-x`},
			options: commonOptions,
			flags: []parserTestCaseFlag{
				{
					flag:  `x`,
					value: true,
				},
			},
		},
	}
	runTests(testCases, t)
}
