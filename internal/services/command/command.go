package command

import (
	"project-s/internal/types"
)

type Command struct {
	CommandList map[string]func(types.PlayerAction)
}

func (c *Command) Registry(commName string, comm func(types.PlayerAction)) {
	c.CommandList[commName] = comm
}

func (c *Command) Use(commName string, playerAction types.PlayerAction) {
	if Use, exists := c.CommandList[commName]; exists {
		Use(playerAction)
	}
}
