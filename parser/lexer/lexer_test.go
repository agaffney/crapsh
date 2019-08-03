package lexer

import (
	"strings"
	"testing"
)

type testCaseOutput struct {
	name  string
	value string
}

type testCase struct {
	input  string
	output []testCaseOutput
}

func TestLexer(t *testing.T) {
	test_cases := []testCase{
		{
			input: "echo foo bar",
			output: []testCaseOutput{
				{name: `Identifier`, value: `echo`},
				{name: `Whitespace`, value: ` `},
				{name: `Identifier`, value: `foo`},
				{name: `Whitespace`, value: ` `},
				{name: `Identifier`, value: `bar`},
			},
		},
		{
			input: "foo $(echo bar foo bar) baz\nabc \"123 456\" 'd\nef' 789",
			output: []testCaseOutput{
				{name: `Identifier`, value: `foo`},
				{name: `Whitespace`, value: ` `},
				{name: `SubshellCaptureOpen`, value: `$(`},
				{name: `Identifier`, value: `echo`},
				{name: `Whitespace`, value: ` `},
				{name: `Identifier`, value: `bar`},
				{name: `Whitespace`, value: ` `},
				{name: `Identifier`, value: `foo`},
				{name: `Whitespace`, value: ` `},
				{name: `Identifier`, value: `bar`},
				{name: `ParenClose`, value: `)`},
				{name: `Whitespace`, value: ` `},
				{name: `Identifier`, value: `baz`},
				{name: `Newline`, value: "\n"},
				{name: `Identifier`, value: `abc`},
				{name: `Whitespace`, value: ` `},
				{name: `DoubleQuote`, value: `"`},
				{name: `Generic`, value: `123`},
				{name: `Whitespace`, value: ` `},
				{name: `Generic`, value: `456`},
				{name: `DoubleQuote`, value: `"`},
				{name: `Whitespace`, value: ` `},
				{name: `SingleQuote`, value: `'`},
				{name: `Identifier`, value: `d`},
				{name: `Newline`, value: "\n"},
				{name: `Identifier`, value: `ef`},
				{name: `SingleQuote`, value: `'`},
				{name: `Whitespace`, value: ` `},
				{name: `Generic`, value: `789`},
			},
		},
	}
	lexer := New()
	for _, test_case := range test_cases {
		lexer.Reset()
		lexer.Start(strings.NewReader(test_case.input))
		for _, expected := range test_case.output {
			token := lexer.ReadToken()
			if token == nil {
				t.Fatalf("Expected token `%s` with value `%s`, got nil", expected.name, expected.value)
			}
			if token.Type != expected.name || token.Value != expected.value {
				t.Fatalf("Expected token `%s` with value `%s`, got token `%s` with value `%s`", expected.name, expected.value, token.Type, token.Value)
			}
		}
		token := lexer.ReadToken()
		if token != nil {
			t.Fatalf("Expected nil token, got token `%s` with value `%s`", token.Type, token.Value)
		}
	}
}
