package main

import (
	"errors"
	"fmt"
)

type Command struct {
	Names       []string
	NArgs       int
	Action      func(params []string) error
	Description string
}

type Arity struct {
	Name  string
	NArgs int
}

var Nop = errors.New("nop")

func (c *Command) String() string {
	return fmt.Sprintf("%s/%d", c.Names, c.NArgs)
}

type Dispatcher struct {
	mapping  map[Arity]*Command
	commands []*Command
}

func NewDispatcher(commands []*Command) *Dispatcher {
	mapping := make(map[Arity]*Command)
	for _, c := range commands {
		for _, name := range c.Names {
			mapping[Arity{Name: name, NArgs: c.NArgs}] = c
		}
	}
	return &Dispatcher{
		mapping:  mapping,
		commands: commands,
	}
}

func (d *Dispatcher) Lookup(blocks []string) (*Command, []string) {
	if len(blocks) == 0 {
		return nil, nil
	}
	com, ok := d.mapping[Arity{Name: blocks[0], NArgs: len(blocks) - 1}]
	if !ok {
		return nil, nil
	}
	return com, blocks[1:]
}

func (d *Dispatcher) Call(line string) error {
	blocks := SplitCommand(line)
	command, params := d.Lookup(blocks)
	if command == nil {
		return Nop
	}
	return command.Action(params)
}
