package pad

import "sync"

type Lighter interface {
	Light(x, y, g, r int) error
	Clear() error
}

type triggerAreaLighter struct {
	page     *Page
	delegate Lighter
}

func (t *triggerAreaLighter) Light(x, y, g, r int) error {
	//make sure that only trigger buttons can be lighted
	if x >= 0 && x < 8 &&
		y >= 0 && y < 8 {

		return t.delegate.Light(x, y, g, r)
	}

	return nil
}

func (t *triggerAreaLighter) Clear() error {
	//clear the whole pad ...
	if err := t.delegate.Clear(); err != nil {
		return err
	}

	//redraw the page lights
	return t.page.Light(t.delegate)
}

type threadSafeLighter struct {
	mux      sync.Mutex
	delegate Lighter
}

func (t *threadSafeLighter) Light(x, y, g, r int) error {
	t.mux.Lock()
	defer t.mux.Unlock()

	return t.delegate.Light(x, y, g, r)
}

func (t *threadSafeLighter) Clear() error {
	t.mux.Lock()
	defer t.mux.Unlock()

	return t.delegate.Clear()
}
