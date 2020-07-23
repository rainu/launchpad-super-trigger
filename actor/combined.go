package actor

import (
	"context"
	"errors"
	"strings"
	"sync"
)

type SequentialActor struct {
	delegates []Actor
}

func (s *SequentialActor) AddActor(actor Actor) {
	s.delegates = append(s.delegates, actor)
}

func (s *SequentialActor) Do(ctx context.Context) error {
	for _, delegate := range s.delegates {
		if err := delegate.Do(ctx); err != nil {
			return err
		}
	}

	return nil
}

type ParallelActor struct {
	delegates []Actor
}

func (p *ParallelActor) AddActor(actor Actor) {
	p.delegates = append(p.delegates, actor)
}

func (p *ParallelActor) Do(ctx context.Context) error {
	wg := sync.WaitGroup{}
	errChan := make(chan error, len(p.delegates))

	doActor := func(actor Actor) {
		defer wg.Done()

		if err := actor.Do(ctx); err != nil {
			errChan <- err
		}
	}

	for _, delegate := range p.delegates {

		wg.Add(1)
		go doActor(delegate)
	}

	wg.Wait()
	close(errChan)

	sb := strings.Builder{}
	isFirstErr := true
	for err := range errChan {
		if isFirstErr {
			sb.WriteString("['")
		} else {
			sb.WriteString("' ; '")
		}

		sb.WriteString(err.Error())

		isFirstErr = false
	}

	if sb.Len() > 0 {
		sb.WriteString("']")
		return errors.New(sb.String())
	}

	return nil
}
