package gui

import (
	"fmt"
	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	realLaunchpad "github.com/rainu/launchpad"
	"go.uber.org/zap"
)

const (
	sizeButtonGridHeight    = 32
	sizeButtonGridWidth     = 32
	sizeButtonSpecialHeight = sizeButtonGridHeight
	sizeButtonSpecialWidth  = sizeButtonGridWidth
)

func (f *FakeLaunchpad) setupWindow() {
	f.window = f.app.NewWindow("Launchpad")
	f.window.Resize(fyne.Size{
		Height: sizeButtonGridHeight * 8,
		Width:  sizeButtonGridWidth * 8,
	})
	f.window.SetFixedSize(true)

	f.triggerButtons = make([][]*widget.Button, 8, 8)
	f.specialTopButtons = make([]*widget.Button, 8, 8)
	f.specialRightButtons = make([]*widget.Button, 8, 8)

	allGridButtons := make([]fyne.CanvasObject, 9*9, 9*9)

	for y := 0; y < 8; y++ {
		f.triggerButtons[y] = make([]*widget.Button, 8, 8)
		for x := 0; x < 8; x++ {
			f.triggerButtons[y][x] = widget.NewButton("", f.generateGridButtonTapHandler(x, y))
			f.triggerButtons[y][x].Resize(fyne.Size{
				Height: sizeButtonGridHeight,
				Width:  sizeButtonGridWidth,
			})
			f.triggerButtons[y][x].SetIcon(getColorIcon(0, 0))
		}
	}
	for x := 0; x < 8; x++ {
		f.specialTopButtons[x] = widget.NewButton("", f.generateSpecialTopButtonTapHandler(x))
		f.specialTopButtons[x].Resize(fyne.Size{
			Height: sizeButtonSpecialHeight,
			Width:  sizeButtonSpecialWidth,
		})
		f.specialTopButtons[x].SetIcon(getColorIcon(0, 0))
	}

	for y := 0; y < 8; y++ {
		f.specialRightButtons[y] = widget.NewButton("", f.generateSpecialRightButtonTapHandler(y))
		f.specialRightButtons[y].Resize(fyne.Size{
			Height: sizeButtonSpecialHeight,
			Width:  sizeButtonSpecialWidth,
		})
		f.specialRightButtons[y].SetIcon(getColorIcon(0, 0))
	}

	i := 0
	for _, button := range f.specialTopButtons {
		allGridButtons[i] = button
		i++
	}

	//right top corner
	allGridButtons[i] = widget.NewLabel("LP")
	i++

	s := 0
	for _, buttonRow := range f.triggerButtons {
		for _, button := range buttonRow {
			allGridButtons[i] = button
			i++
		}

		allGridButtons[i] = f.specialRightButtons[s]

		s++
		i++
	}

	f.window.SetContent(widget.NewVBox(
		fyne.NewContainerWithLayout(
			layout.NewGridLayout(9),
			allGridButtons...,
		),
	))

	f.window.SetOnClosed(func() {
		f.isOpen = false
	})
}

func (f *FakeLaunchpad) generateGridButtonTapHandler(x, y int) func() {
	return func() {
		zap.L().Info(fmt.Sprintf("TAP: %d, %d", x, y))
		f.hitChannel <- realLaunchpad.Hit{
			X:    x,
			Y:    y,
			Down: true,
		}
		f.hitChannel <- realLaunchpad.Hit{
			X:    x,
			Y:    y,
			Down: false,
		}
	}
}

func (f *FakeLaunchpad) generateSpecialTopButtonTapHandler(x int) func() {
	return func() {
		zap.L().Info(fmt.Sprintf("TAP Special Top: %d", x))
		f.hitChannel <- realLaunchpad.Hit{
			X:    x,
			Y:    8,
			Down: true,
		}
		f.hitChannel <- realLaunchpad.Hit{
			X:    x,
			Y:    8,
			Down: false,
		}
	}
}

func (f *FakeLaunchpad) generateSpecialRightButtonTapHandler(y int) func() {
	return func() {
		zap.L().Info(fmt.Sprintf("TAP Special Right: %d", y))
		f.hitChannel <- realLaunchpad.Hit{
			X:    8,
			Y:    y,
			Down: true,
		}
		f.hitChannel <- realLaunchpad.Hit{
			X:    8,
			Y:    y,
			Down: false,
		}
	}
}
