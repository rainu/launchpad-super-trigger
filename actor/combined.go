package actor

import (
	"errors"
	"strings"
	"sync"
)

type MetaActor interface {
	HasActor(func(Actor) bool) bool
}

type Sequential struct {
	delegates []Actor
}

func (s *Sequential) AddActor(actor Actor) {
	s.delegates = append(s.delegates, actor)
}

func (s *Sequential) HasActor(clb func(Actor) bool) bool {
	for _, delegate := range s.delegates {
		if clb(delegate) {
			return true
		}
	}

	return false
}

func (s *Sequential) Do(ctx Context) error {
	for _, delegate := range s.delegates {
		if err := delegate.Do(ctx); err != nil {
			return err
		}
	}

	return nil
}

type Parallel struct {
	delegates []Actor
}

func (p *Parallel) AddActor(actor Actor) {
	p.delegates = append(p.delegates, actor)
}

func (p *Parallel) HasActor(clb func(Actor) bool) bool {
	for _, delegate := range p.delegates {
		if clb(delegate) {
			return true
		}
	}

	return false
}

func (p *Parallel) Do(ctx Context) error {
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
