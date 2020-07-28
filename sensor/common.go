package sensor

import (
	"context"
	"sync"
)

type Callback func(Sensor)

type Sensor interface {
	Run(ctx context.Context) error
	AddCallback(clbId interface{}, clb Callback)
	RemoveCallback(clbId interface{})
	LastMessage() []byte
}

type callbackHandler struct {
	callbacks     map[interface{}]Callback
	callbackMutex sync.RWMutex
}

func (c *callbackHandler) AddCallback(clbId interface{}, clb Callback) {
	c.callbackMutex.Lock()
	defer c.callbackMutex.Unlock()

	if c.callbacks == nil {
		c.callbacks = map[interface{}]Callback{}
	}

	c.callbacks[clbId] = clb
}

func (c *callbackHandler) RemoveCallback(clbId interface{}) {
	c.callbackMutex.Lock()
	defer c.callbackMutex.Unlock()

	if c.callbacks == nil {
		return
	}

	delete(c.callbacks, clbId)
}

func (c *callbackHandler) Call(sensor Sensor) {
	c.callbackMutex.RLock()
	defer c.callbackMutex.RUnlock()

	for _, callback := range c.callbacks {
		callback(sensor)
	}
}
