package parser

import (
	"github.com/agaffney/crapsh/parser/tokens"
)

type Container struct {
	Name              string
	Token             string
	TokenEnd          string
	Escapes           bool
	AllowedContainers []string
}

var containers []*Container
var line_container *Container

func init() {
	line_container = &Container{
		Name:              `Line`,
		TokenEnd:          tokens.NEWLINE,
		Escapes:           true,
		AllowedContainers: []string{"StringSingle", "StringDouble", "Variable"},
	}
	containers = []*Container{
		{
			Name:     `StringSingle`,
			Token:    tokens.SINGLE_QUOTE,
			TokenEnd: tokens.SINGLE_QUOTE,
			Escapes:  false,
		},
		{
			Name:              `StringDouble`,
			Token:             tokens.DOUBLE_QUOTE,
			TokenEnd:          tokens.DOUBLE_QUOTE,
			Escapes:           true,
			AllowedContainers: []string{"Variable", "Subshell"},
		},
		{
			Name:     `Variable`,
			Token:    tokens.VARIABLE_OPEN,
			TokenEnd: tokens.CURLY_BRACE_CLOSE,
			Escapes:  false,
		},
	}
}

func (c *Container) Allowed_container(s string) bool {
	for _, foo := range c.AllowedContainers {
		if s == foo {
			return true
		}
	}
	return false
}
