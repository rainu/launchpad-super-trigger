package sensor

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"os/exec"
	"sync"
	"time"
)

type Command struct {
	callbackHandler

	Name      string
	Arguments []string
	Interval  time.Duration

	running     bool
	mux         sync.RWMutex
	lastMessage []byte
}

func (c *Command) Run(ctx context.Context) error {
	if c.running {
		return fmt.Errorf("listerner is already running")
	}

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

	c.mux.Lock()
	c.lastMessage = out
	c.mux.Unlock()

	c.callbackHandler.Call(c)
	return execErr
}

func (c *Command) LastMessage() []byte {
	c.mux.RLock()
	defer c.mux.RUnlock()

	return c.lastMessage
}
