package ast

import ()

type Command struct {
	NodeBase
}

func NewCommand() Node {
	c := &Command{NodeBase: NodeBase{Name: `Command`}}
	return c
}
