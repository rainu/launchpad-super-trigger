package sensor

import (
	"context"
	"github.com/rainu/launchpad-super-trigger/sensor/store"
	"go.uber.org/zap"
)

type Static struct {
	callbackHandler

	MessageStore store.Store

	running bool
}

func (s *Static) Run(ctx context.Context) error {
	//wait until context closed
	<-ctx.Done()

	return nil
}

func (s *Static) Set(state []byte) {
	if err := s.MessageStore.Set(state); err != nil {
		zap.L().Error("Could not save message into message store!", zap.Error(err))
	}
	s.callbackHandler.Call(s)
}

func (s *Static) LastMessage() []byte {
	data, err := s.MessageStore.Get()
	if err != nil {
		zap.L().Error("Could not get message from message store!", zap.Error(err))
		return nil
	}

	return data
}
