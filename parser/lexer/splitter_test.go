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

func runTests(testCases []splitterTestCase, t *testing.T) {
	lexer := New()
	for _, test_case := range testCases {
		lexer.Reset()
		input := parser_input.NewStringParserInput(test_case.input)
		lexer.Start(input)
		for idx, expected := range test_case.output {
			token, err := lexer.ReadToken()
			if err != nil {
				if err == io.EOF {
					// Restart lexer to read another line
					lexer.Start(input)
					// Read a new token if we didn't get one before
					if token == nil {
						var err2 error
						token, err2 = lexer.ReadToken()
						if err2 != nil {
							if err2 == io.EOF {
								if idx < len(test_case.output)-1 {
									t.Fatalf("Encountered unexpected EOF")
								}
							} else {
								t.Fatalf("Encountered unexpected error: %s", err2.Error())
							}
						}
					}
				} else {
					t.Fatalf("Encountered unexpected error: %s", err.Error())
				}
			}
			if token == nil {
				t.Fatalf("Expected token type %d with value `%s`, got nil", expected.token, expected.value)
			}
			if token.Type != expected.token || token.Value != expected.value {
				t.Fatalf("Expected token type %d with value `%s`, got token type %d with value `%s`", expected.token, expected.value, token.Type, token.Value)
			}
		}
		token, err := lexer.ReadToken()
		if err != nil && err != io.EOF {
			t.Fatalf("Unexpected error: %s", err.Error())
		}
		if token != nil {
			t.Fatalf("Expected nil token, got token type %d with value `%s`", token.Type, token.Value)
		}
	}
}

func TestSplitterBasic(t *testing.T) {
	test_cases := []splitterTestCase{
		{
			input: "echo foo bar",
			output: []splitterTestCaseOutput{
				{token: tokens.TOKEN_NULL, value: `echo`},
				{token: tokens.TOKEN_NULL, value: `foo`},
				{token: tokens.TOKEN_NULL, value: `bar`},
			},
		},
	}
	runTests(test_cases, t)
}

func TestSplitterMultipleLines(t *testing.T) {
	test_cases := []splitterTestCase{
		{
			input: "foo $(echo bar foo bar) baz\nabc \"123 456\" 789",
			output: []splitterTestCaseOutput{
				{token: tokens.TOKEN_NULL, value: `foo`},
				{token: tokens.TOKEN_NULL, value: `$(echo bar foo bar)`},
				{token: tokens.TOKEN_NULL, value: `baz`},
				{token: tokens.TOKEN_NEWLINE, value: "\n"},
				{token: tokens.TOKEN_NULL, value: `abc`},
				{token: tokens.TOKEN_NULL, value: `"123 456"`},
				{token: tokens.TOKEN_NULL, value: `789`},
			},
		},
	}
	runTests(test_cases, t)
}

func TestSplitterComments(t *testing.T) {
	test_cases := []splitterTestCase{
		{
			input: "echo foo bar#baz\necho foo #bar",
			output: []splitterTestCaseOutput{
				{token: tokens.TOKEN_NULL, value: `echo`},
				{token: tokens.TOKEN_NULL, value: `foo`},
				{token: tokens.TOKEN_NULL, value: `bar`},
				{token: tokens.TOKEN_NEWLINE, value: "\n"},
				{token: tokens.TOKEN_NULL, value: `echo`},
				{token: tokens.TOKEN_NULL, value: `foo`},
			},
		},
	}
	runTests(test_cases, t)
}

func TestSplitterContinuation(t *testing.T) {
	test_cases := []splitterTestCase{
		{
			input: "echo foo 'bar\nbaz' abc 123",
			output: []splitterTestCaseOutput{
				{token: tokens.TOKEN_NULL, value: `echo`},
				{token: tokens.TOKEN_NULL, value: `foo`},
				{token: tokens.TOKEN_NULL, value: "'bar\nbaz'"},
				{token: tokens.TOKEN_NULL, value: `abc`},
				{token: tokens.TOKEN_NULL, value: `123`},
			},
		},
		{
			input: "echo foo bar \\\nbaz abc 123",
			output: []splitterTestCaseOutput{
				{token: tokens.TOKEN_NULL, value: `echo`},
				{token: tokens.TOKEN_NULL, value: `foo`},
				{token: tokens.TOKEN_NULL, value: `bar`},
				{token: tokens.TOKEN_NULL, value: `baz`},
				{token: tokens.TOKEN_NULL, value: `abc`},
				{token: tokens.TOKEN_NULL, value: `123`},
			},
		},
	}
	runTests(test_cases, t)
}
