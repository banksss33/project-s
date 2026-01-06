package command

import (
	"project-s/internal/types"
)

type Command struct {
	CommandList map[string]func(types.ActionEvent)
}

func (c *Command) Registry(commName string, comm func(types.ActionEvent)) {
	c.CommandList[commName] = comm
}

func (c *Command) Use(commName string, jsonData types.ActionEvent) {
	if GetCommand, exists := c.CommandList[commName]; exists {
		GetCommand(jsonData)
	}
}
