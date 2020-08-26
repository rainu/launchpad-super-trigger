package pad

import (
	"github.com/rainu/launchpad"
	"gitlab.com/gomidi/midi"
	"io"
	"strings"
)

type Launchpad interface {
	Lighter
	io.Closer

	ListenToHits() (<-chan launchpad.Hit, error)
	IsHealthy() bool
}

type realLaunchpad struct {
	launchpad.Launchpad

	driver midi.Driver
}

func NewLaunchpad(driver midi.Driver) (Launchpad, error) {
	pad, err := launchpad.NewLaunchpad(driver)
	if err != nil {
		return nil, err
	}
	return &realLaunchpad{
		Launchpad: pad,
		driver:    driver,
	}, nil
}

func (r *realLaunchpad) IsHealthy() bool {
	ins, err := r.driver.Ins()
	if err != nil {
		return false
	}

	for i := range ins {
		if strings.Contains(ins[i].String(), "Launchpad") {
			return true
		}
	}
	return false //no launchpad found -> connections lost
}
