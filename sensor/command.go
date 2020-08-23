package sensor

import (
	"context"
	"fmt"
	"github.com/rainu/launchpad-super-trigger/sensor/store"
	"go.uber.org/zap"
	"os/exec"
	"time"
)

type Command struct {
	callbackHandler

	Name         string
	Arguments    []string
	Interval     time.Duration
	MessageStore store.Store

	running bool
}

func (c *Command) Run(ctx context.Context) error {
	if c.running {
		return fmt.Errorf("listerner is already running")
	}
	defer func() {
		c.running = false
	}()

	timer := time.NewTimer(c.Interval)

	c.running = true

	//first call
	if err := c.call(ctx); err != nil {
		zap.L().Error("Error while call command sensor!", zap.Error(err))
	}

	for {
		select {
		case <-timer.C:
			if err := c.call(ctx); err != nil {
				zap.L().Error("Error while call command sensor!", zap.Error(err))
			}

			//reset timer
			timer = time.NewTimer(c.Interval)

		//wait until context closed
		case <-ctx.Done():
			return nil
		}
	}
}

func (c *Command) call(ctx context.Context) error {
	command := exec.CommandContext(ctx, c.Name, c.Arguments...)

	out, execErr := command.CombinedOutput()
	zap.L().Info("Command call was successful!")

	if err := c.MessageStore.Set(out); err != nil {
		zap.L().Error("Could not save message into message store!", zap.Error(err))
	}

	c.callbackHandler.Call(c)
	return execErr
}

func (c *Command) LastMessage() []byte {
	data, err := c.MessageStore.Get()
	if err != nil {
		zap.L().Error("Could not get message from message store!", zap.Error(err))
		return nil
	}

	return data
}
