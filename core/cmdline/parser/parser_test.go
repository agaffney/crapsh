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
		args, err := Parse(testCase.options)
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
		// TODO: check flags
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
					flag: `x`,
				},
			},
		},
	}
	runTests(testCases, t)
}
