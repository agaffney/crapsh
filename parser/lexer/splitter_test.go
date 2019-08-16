package lexer

import (
	"bufio"
	"github.com/agaffney/crapsh/parser/tokens"
	"io"
	"strings"
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

type parserInput struct {
	input *bufio.Reader
}

func NewParserInput(input string) *parserInput {
	i := &parserInput{}
	i.input = bufio.NewReader(strings.NewReader(input))
	return i
}

func (i *parserInput) ReadLine() (string, error) {
	return i.input.ReadString('\n')
}

func (i *parserInput) ReadAnotherLine() (string, error) {
	return i.input.ReadString('\n')
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
				{token: tokens.TOKEN_NULL, value: `'d\nef'`},
				{token: tokens.TOKEN_NULL, value: `789`},
			},
		},
	}
	lexer := New()
	for _, test_case := range test_cases {
		lexer.Reset()
		input := NewParserInput(test_case.input)
		lexer.Start(input)
		for _, expected := range test_case.output {
			token, err := lexer.NextToken()
			if err != nil {
				if err == io.EOF {
					break
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
