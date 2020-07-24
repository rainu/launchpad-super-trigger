package pad

import (
	"github.com/rakyll/launchpad"
	"github.com/rakyll/portmidi"
	"reflect"
	"unsafe"
)

type ScrollingTextBuilder interface {
	Add(speed byte, text string) ScrollingTextBuilder
	Perform() error
}

type scrollingText struct {
	seq        []byte
	midiStream **portmidi.Stream
}

type ExtendedLaunchpad interface {
	Listen() <-chan launchpad.Hit
	Read() (hits []launchpad.Hit, err error)
	Light(x int, y int, g int, r int) error
	Clear() error
	Close() error

	//Text will return a scrolling text builder whether you can build and
	//execute an text which will be scrolled on the launchpad
	Text(g int, r int) ScrollingTextBuilder

	//TextLoop will return a scrolling text builder whether you can build and
	//execute an text which will be scrolled endless on the launchpad. If you
	//want to stop an text loop you have to build and execute an empty textLoop!
	TextLoop(g int, r int) ScrollingTextBuilder
}

type extendedLaunchpad struct {
	launchpad  *launchpad.Launchpad
	midiStream **portmidi.Stream
}

func newExtendedLaunchpad(launchpad *launchpad.Launchpad) ExtendedLaunchpad {
	pointerVal := reflect.ValueOf(launchpad)
	val := reflect.Indirect(pointerVal)

	member := val.FieldByName("outputStream")
	ptrToY := unsafe.Pointer(member.UnsafeAddr())
	realPtrToY := (**portmidi.Stream)(ptrToY)

	return &extendedLaunchpad{
		launchpad:  launchpad,
		midiStream: realPtrToY,
	}
}

func (e extendedLaunchpad) Listen() <-chan launchpad.Hit {
	return e.launchpad.Listen()
}

func (e extendedLaunchpad) Read() (hits []launchpad.Hit, err error) {
	return e.launchpad.Read()
}

func (e extendedLaunchpad) Light(x int, y int, g int, r int) error {
	return e.launchpad.Light(x, y, g, r)
}

func (e extendedLaunchpad) Clear() error {
	return e.launchpad.Clear()
}

func (e extendedLaunchpad) Close() error {
	return e.launchpad.Close()
}

func (e extendedLaunchpad) Text(g int, r int) ScrollingTextBuilder {
	color := byte(16*g + r + 8 + 4)

	return &scrollingText{
		seq:        []byte{0xF0, 0x00, 0x20, 0x29, 0x09, color, 0x05},
		midiStream: e.midiStream,
	}
}

func (e extendedLaunchpad) TextLoop(g int, r int) ScrollingTextBuilder {
	color := byte(16*g + r + 8 + 4 + 64)

	return &scrollingText{
		seq:        []byte{0xF0, 0x00, 0x20, 0x29, 0x09, color},
		midiStream: e.midiStream,
	}
}

func (s *scrollingText) Add(speed byte, text string) ScrollingTextBuilder {
	if speed > 7 {
		speed = 7
	} else if speed < 1 {
		speed = 1
	}

	s.seq = append(s.seq, speed)
	s.seq = append(s.seq, []byte(text)...)

	return s
}

func (s *scrollingText) Perform() error {
	s.seq = append(s.seq, 0xF7)

	return (*s.midiStream).WriteSysExBytes(portmidi.Time(), s.seq)
}
