package actor

import (
	"fmt"
	"os/exec"
)

type Command struct {
	Name          string
	Arguments     []string
	AppendContext bool
}

func (c *Command) Do(ctx Context) error {
	var args []string
	args = c.Arguments

	if c.AppendContext {
		args = append(args,
			fmt.Sprintf("%d", ctx.Page),
			fmt.Sprintf("%d", ctx.HitX),
			fmt.Sprintf("%d", ctx.HitY),
		)
	}

	command := exec.CommandContext(ctx.Context, c.Name, args...)

	return command.Run()
}
