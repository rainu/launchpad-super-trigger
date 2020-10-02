package pad

const (
	PageNavigationBinary = byte(0)
	PageNavigationToggle = byte(1)
)

type special struct {
	pageNavigationMode byte
	locked             bool
}

func newSpecial() *special {
	return &special{
		pageNavigationMode: PageNavigationBinary,
		locked:             false,
	}
}

// Light will apply the color representations of special buttons (right row)
func (s *special) Light(pad Lighter) error {
	if err := s.lightNavigation(pad); err != nil {
		return err
	}

	if err := s.lightArm(pad); err != nil {
		return err
	}

	return nil
}

func (s *special) SetPageNavigationMode(mode byte, pad Lighter) error {
	s.pageNavigationMode = mode
	return s.Light(pad)
}

func (s *special) SwitchPageNavigationMode(pad Lighter) error {
	if s.pageNavigationMode == PageNavigationBinary {
		s.pageNavigationMode = PageNavigationToggle
	} else {
		s.pageNavigationMode = PageNavigationBinary
	}

	return s.Light(pad)
}

func (s *special) lightNavigation(pad Lighter) error {
	if s.pageNavigationMode == PageNavigationBinary {
		return ColorOff.Light(pad, 8, 0)
	} else {
		return ColorGreen.Light(pad, 8, 0)
	}
}

func (s *special) ToggleLock(pad Lighter) error {
	s.locked = !s.locked

	return s.lightArm(pad)
}

func (s *special) lightArm(pad Lighter) error {
	if s.locked {
		return ColorRed.Light(pad, 8, 7)
	}
	return ColorOff.Light(pad, 8, 7)
}
