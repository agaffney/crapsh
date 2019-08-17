package lexer

import (
	parser_input "github.com/agaffney/crapsh/parser/input"
	"github.com/agaffney/crapsh/parser/tokens"
	"io"
	"testing"
)

type splitterTestCaseOutput struct {
	token int
	value string
}

type splitterTestCase struct {
	input  string
	output []splitterTestCaseOutput
}

func TestSplitter(t *testing.T) {
	test_cases := []splitterTestCase{
		{
			input: "echo foo bar",
			output: []splitterTestCaseOutput{
				{token: tokens.TOKEN_NULL, value: `echo`},
				{token: tokens.TOKEN_NULL, value: `foo`},
				{token: tokens.TOKEN_NULL, value: `bar`},
			},
		},
		{
			input: "foo $(echo bar foo bar) baz\nabc \"123 456\" 'd\nef' 789",
			output: []splitterTestCaseOutput{
				{token: tokens.TOKEN_NULL, value: `foo`},
				{token: tokens.TOKEN_NULL, value: `$(echo bar foo bar)`},
				{token: tokens.TOKEN_NULL, value: `baz`},
				{token: tokens.TOKEN_NEWLINE, value: "\n"},
				{token: tokens.TOKEN_NULL, value: `abc`},
				{token: tokens.TOKEN_NULL, value: `"123 456"`},
				{token: tokens.TOKEN_NULL, value: "'d\nef'"},
				{token: tokens.TOKEN_NULL, value: `789`},
			},
		},
	}
	lexer := New()
	for _, test_case := range test_cases {
		lexer.Reset()
		input := parser_input.NewStringParserInput(test_case.input)
		lexer.Start(input)
		for idx, expected := range test_case.output {
			token, err := lexer.NextToken()
			if err != nil {
				if err == io.EOF {
					// Read additional line
					err2 := lexer.readLine(false)
					if err2 != nil {
						if err2 == io.EOF {
							if idx < len(test_case.output)-1 {
								t.Fatalf("Encountered unexpected EOF")
							}
						} else {
							t.Fatalf("Encountered unexpected error: %s", err2.Error())
						}
					}
					// Read a new token if we didn't get one before
					if token == nil {
						token, _ = lexer.NextToken()
					}
				} else {
					t.Fatalf("Unexpected error: %s", err.Error())
				}
			}
			if token == nil {
				t.Fatalf("Expected token type %d with value `%s`, got nil", expected.token, expected.value)
			}
			if token.Type != expected.token || token.Value != expected.value {
				t.Fatalf("Expected token type %d with value `%s`, got token type %d with value `%s`", expected.token, expected.value, token.Type, token.Value)
			}
		}
		token, err := lexer.NextToken()
		if err != nil && err != io.EOF {
			t.Fatalf("Unexpected error: %s", err.Error())
		}
		if token != nil {
			t.Fatalf("Expected nil token, got token type %d with value `%s`", token.Type, token.Value)
		}
	}
}
