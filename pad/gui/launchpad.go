package gui

import (
	"fmt"
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
	realLaunchpad "github.com/rainu/launchpad"
)

type FakeLaunchpad struct {
	app    fyne.App
	window fyne.Window

	triggerButtons      [][]*widget.Button
	specialTopButtons   []*widget.Button
	specialRightButtons []*widget.Button
	hitChannel          chan realLaunchpad.Hit

	isOpen bool
}

func NewLaunchpad() *FakeLaunchpad {
	result := &FakeLaunchpad{
		isOpen:     true,
		hitChannel: make(chan realLaunchpad.Hit, 9*9),
	}

	result.app = app.New()
	result.setupWindow()

	return result
}

func (f FakeLaunchpad) StartAndBlock() {
	f.window.ShowAndRun()
}

func (f FakeLaunchpad) Light(x, y, g, r int) error {
	if x >= 0 && x <= 7 && y >= 0 && y <= 7 {
		f.triggerButtons[y][x].SetIcon(getColorIcon(r, g))
	} else if x == 8 {
		f.specialRightButtons[y].SetIcon(getColorIcon(r, g))
	} else if y == 8 {
		f.specialTopButtons[x].SetIcon(getColorIcon(r, g))
	} else {
		return fmt.Errorf("invalid coordinate (out of range)")
	}

	return nil
}

func (f FakeLaunchpad) Text(g, r int) realLaunchpad.ScrollingTextBuilder {
	return nil
}

func (f FakeLaunchpad) TextLoop(g, r int) realLaunchpad.ScrollingTextBuilder {
	return nil
}

func (f FakeLaunchpad) Clear() error {
	for _, button := range f.specialTopButtons {
		button.SetIcon(getColorIcon(0, 0))
	}
	for _, buttonRow := range f.triggerButtons {
		for _, button := range buttonRow {
			button.SetIcon(getColorIcon(0, 0))
		}
	}
	for _, button := range f.specialRightButtons {
		button.SetIcon(getColorIcon(0, 0))
	}
	return nil
}

func (f FakeLaunchpad) Close() error {
	return nil
}

func (f FakeLaunchpad) ListenToHits() (<-chan realLaunchpad.Hit, error) {
	return f.hitChannel, nil
}

func (f FakeLaunchpad) IsHealthy() bool {
	return f.isOpen
}
