package pad

import (
	"fmt"
	"go.uber.org/zap"
)

type handlerId struct {
	Number PageNumber
	X      int
	Y      int
}

type Handler interface {
	// OnTrigger will called each time a hit was made.
	OnTrigger(Lighter, *LaunchpadSuperTrigger, PageNumber, int, int) error

	// OnPageEnter will called each time a page will be entered.
	OnPageEnter(Lighter, *LaunchpadSuperTrigger, PageNumber) error

	// OnPageLeave will called each time a page will be leaved.
	OnPageLeave(Lighter, *LaunchpadSuperTrigger, PageNumber) error
}

type PageHandleFunc func(lighter Lighter, page PageNumber) error

type SimpleHandler struct {
	TriggerFn   TriggerHandleFunc
	PageEnterFn PageHandleFunc
	PageLeaveFn PageHandleFunc
}

type TriggerDispatcher struct {
	currentPage      PageNumber
	triggerDelegates map[string]Handler
	pageDelegates    map[PageNumber]Handler
	firstEnterCalled bool
}

func (t *TriggerDispatcher) Handle(lighter Lighter, lst *LaunchpadSuperTrigger, page PageNumber, x, y int) error {
	if !t.firstEnterCalled {
		if handler := t.lookupPageHandler(page); handler != nil {
			if err := handler.OnPageEnter(lighter, lst, page); err != nil {
				return err
			}
		}
		t.firstEnterCalled = true
	}

	if t.currentPage != page {
		lastPage := t.currentPage
		t.currentPage = page

		//change page
		if handler := t.lookupPageHandler(lastPage); handler != nil {
			if err := handler.OnPageLeave(lighter, lst, lastPage); err != nil {
				return err
			}
		}
		if err := lighter.Clear(); err != nil {
			return err
		}
		if handler := t.lookupPageHandler(page); handler != nil {
			if err := handler.OnPageEnter(lighter, lst, page); err != nil {
				return err
			}
		}
	}

	if x < 0 || y < 0 {
		//skip hits outside the trigger area
		return nil
	}

	if handler := t.lookupTriggerHandler(page, x, y); handler != nil {
		zap.L().Debug(fmt.Sprintf("Call handler for page %d %d:%d", page, x, y))
		return handler.OnTrigger(lighter, lst, page, x, y)
	}
	zap.L().Warn(fmt.Sprintf("No handler found for page %d %d:%d", page, x, y))

	return nil
}

func (t *TriggerDispatcher) AddPageHandler(handler Handler, page ...PageNumber) {
	for _, number := range page {
		for x := 0; x < 8; x++ {
			for y := 0; y < 8; y++ {
				t.AddSpecificHandler(handler, number, x, y)
			}
		}
	}
}

func (t *TriggerDispatcher) AddSpecificHandler(handler Handler, page PageNumber, x, y int) {
	if t.triggerDelegates == nil {
		t.triggerDelegates = map[string]Handler{}
	}
	if t.pageDelegates == nil {
		t.pageDelegates = map[PageNumber]Handler{}
	}

	t.triggerDelegates[handlerId{
		Number: page,
		X:      x,
		Y:      y,
	}.Id()] = handler
	t.pageDelegates[page] = handler
}

func (t *TriggerDispatcher) lookupTriggerHandler(number PageNumber, x, y int) Handler {
	return t.triggerDelegates[handlerId{
		Number: number,
		X:      x,
		Y:      y,
	}.Id()]
}

func (t *TriggerDispatcher) lookupPageHandler(number PageNumber) Handler {
	return t.pageDelegates[number]
}

func (h handlerId) Id() string {
	return fmt.Sprintf("%d_%d_%d", h.Number, h.X, h.Y)
}

func (s *SimpleHandler) OnTrigger(lighter Lighter, lst *LaunchpadSuperTrigger, page PageNumber, x, y int) error {
	if s.TriggerFn != nil {
		return s.TriggerFn(lighter, lst, page, x, y)
	}
	return nil
}

func (s *SimpleHandler) OnPageEnter(lighter Lighter, lst *LaunchpadSuperTrigger, page PageNumber) error {
	if s.PageEnterFn != nil {
		return s.PageEnterFn(lighter, page)
	}
	return nil
}

func (s *SimpleHandler) OnPageLeave(lighter Lighter, lst *LaunchpadSuperTrigger, page PageNumber) error {
	if s.PageLeaveFn != nil {
		return s.PageLeaveFn(lighter, page)
	}
	return nil
}
