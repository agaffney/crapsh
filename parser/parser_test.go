package parser

import (
	"strings"
	"testing"
)

func TestParser(t *testing.T) {
	parser := NewParser()
	test_inputs := []string{
		"foo $(echo bar foo bar) baz\nabc \"123 456\" 'd\nef' 789",
	}
	for _, input := range test_inputs {
		parser.Parse(strings.NewReader(input))
	}
}
