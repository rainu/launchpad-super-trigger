package pad

const (
	pageNavigationBinary = byte(0)
	pageNavigationToggle = byte(1)
)

type special struct {
	pageNavigationMode byte
	locked             bool
}

func newSpecial() *special {
	return &special{
		pageNavigationMode: pageNavigationBinary,
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

func (s *special) SwitchPageNavigationMode(pad Lighter) error {
	if s.pageNavigationMode == pageNavigationBinary {
		s.pageNavigationMode = pageNavigationToggle
	} else {
		s.pageNavigationMode = pageNavigationBinary
	}

	return s.Light(pad)
}

func (s *special) lightNavigation(pad Lighter) error {
	if s.pageNavigationMode == pageNavigationBinary {
		return ColorOff.Light(pad, 8, 0)
	} else {
		return ColorHighGreen.Light(pad, 8, 0)
	}
}

func (s *special) ToggleLock(pad Lighter) error {
	s.locked = !s.locked

	return s.lightArm(pad)
}

func (s *special) lightArm(pad Lighter) error {
	if s.locked {
		return ColorHighRed.Light(pad, 8, 7)
	}
	return ColorOff.Light(pad, 8, 7)
}
