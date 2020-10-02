package pad

import (
	"github.com/rainu/launchpad"
	"sync"
)

type Lighter interface {
	Light(x, y int, color launchpad.Color) error
	Text(color launchpad.Color) launchpad.ScrollingTextBuilder
	TextLoop(color launchpad.Color) launchpad.ScrollingTextBuilder
	Name() string
	Clear() error
}

type triggerAreaLighter struct {
	page     *Page
	special  *special
	delegate Lighter
}

func (t *triggerAreaLighter) Light(x, y int, color launchpad.Color) error {
	//make sure that only trigger buttons can be lighted
	if x >= 0 && x < 8 &&
		y >= 0 && y < 8 {

		return t.delegate.Light(x, y, color)
	}

	return nil
}

func (t *triggerAreaLighter) Text(color launchpad.Color) launchpad.ScrollingTextBuilder {
	return t.delegate.Text(color)
}

func (t *triggerAreaLighter) TextLoop(color launchpad.Color) launchpad.ScrollingTextBuilder {
	return t.delegate.TextLoop(color)
}

func (t *triggerAreaLighter) Name() string {
	return t.delegate.Name()
}

func (t *triggerAreaLighter) Clear() error {
	//clear the whole pad ...
	if err := t.delegate.Clear(); err != nil {
		return err
	}

	//redraw the page lights
	if err := t.page.Light(t.delegate); err != nil {
		return err
	}

	//...and special lights
	return t.special.Light(t.delegate)
}

type threadSafeLighter struct {
	mux      sync.Mutex
	delegate Lighter
}

type threadSafeTextBuilder struct {
	mux      *sync.Mutex
	delegate launchpad.ScrollingTextBuilder
}

func (t *threadSafeLighter) Light(x, y int, color launchpad.Color) error {
	t.mux.Lock()
	defer t.mux.Unlock()

	return t.delegate.Light(x, y, color)
}

func (t *threadSafeLighter) Text(color launchpad.Color) launchpad.ScrollingTextBuilder {
	return &threadSafeTextBuilder{
		mux:      &t.mux,
		delegate: t.delegate.Text(color),
	}
}

func (t *threadSafeLighter) TextLoop(color launchpad.Color) launchpad.ScrollingTextBuilder {
	return &threadSafeTextBuilder{
		mux:      &t.mux,
		delegate: t.delegate.TextLoop(color),
	}
}

func (t *threadSafeLighter) Name() string {
	return t.delegate.Name()
}

func (t *threadSafeLighter) Clear() error {
	t.mux.Lock()
	defer t.mux.Unlock()

	return t.delegate.Clear()
}

func (t *threadSafeTextBuilder) Add(speed byte, text string) launchpad.ScrollingTextBuilder {
	return t.delegate.Add(speed, text)
}

func (t *threadSafeTextBuilder) Perform() error {
	t.mux.Lock()
	defer t.mux.Unlock()

	return t.delegate.Perform()
}
