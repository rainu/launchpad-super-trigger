package pad

import "sync"

type Lighter interface {
	Light(x, y, g, r int) error
}

type ThreadSafeLighter struct {
	mux      sync.Mutex
	delegate Lighter
}

func (t *ThreadSafeLighter) Light(x, y, g, r int) error {
	t.mux.Lock()
	defer t.mux.Unlock()

	return t.delegate.Light(x, y, g, r)
}
