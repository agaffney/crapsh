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
				{
					name:  `Identifier`,
					value: `echo`,
				},
				{
					name:  `Whitespace`,
					value: ` `,
				},
				{
					name:  `Identifier`,
					value: `foo`,
				},
				{
					name:  `Whitespace`,
					value: ` `,
				},
				{
					name:  `Identifier`,
					value: `bar`,
				},
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
	}
}
