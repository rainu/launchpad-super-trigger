package pad

var setColor = NormalGreen
var unsetColor = Off
var EmptyPage = Page{
	ColorRepresentation: []Color{unsetColor, unsetColor, unsetColor, unsetColor, unsetColor, unsetColor, unsetColor, unsetColor},
}

type Page struct {
	ColorRepresentation []Color
}

func NewPage(number byte) *Page {
	page := Page{
		ColorRepresentation: make([]Color, 8, 8),
	}

	for i := 0; i < 8; i++ {
		c := Off

		if hasBit(number, i) {
			c = setColor
		}

		page.ColorRepresentation[i] = c
	}

	return &page
}

// Light will apply the color representation of this page on the launchpad
func (p *Page) Light(pad Lighter) error {
	if p != &EmptyPage {
		if err := EmptyPage.Light(pad); err != nil {
			return err
		}
	}

	for i := 0; i < 8; i++ {
		if err := p.ColorRepresentation[i].Light(pad, i, 8); err != nil {
			return err
		}
	}

	return nil
}

// Number returns the numeric representation of this page
func (p *Page) Number() byte {
	var n byte
	n = 0

	for i, color := range p.ColorRepresentation {
		if color.Ordinal() == setColor.Ordinal() {
			n = setBit(n, i)
		}
	}

	return n
}

func (p *Page) Toggle(pos int) {
	if p.ColorRepresentation[pos].Ordinal() == setColor.Ordinal() {
		p.ColorRepresentation[pos] = unsetColor
	} else {
		p.ColorRepresentation[pos] = setColor
	}
}

func setBit(n byte, pos int) byte {
	n |= 1 << pos
	return n
}

func hasBit(n byte, pos int) bool {
	val := n & (1 << pos)
	return val > 0
}
