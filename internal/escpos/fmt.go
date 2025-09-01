package escpos

const (
	CharEscape       = 0x1B
	CharBold         = 0x45
	CharDoubleStrike = 0x47
	CharUnderline    = 0x2D
)

func (p *Printer) Emphasize(on bool) {
	if on {
		p.Write([]byte{CharEscape, CharBold, 1})
	} else {
		p.Write([]byte{CharEscape, CharBold, 0})
	}
}

func (p *Printer) DoubleStrike(on bool) {
	if on {
		p.Write([]byte{CharEscape, CharDoubleStrike, 1})
	} else {
		p.Write([]byte{CharEscape, CharDoubleStrike, 0})
	}
}

func (p *Printer) Underline(on bool) {
	if on {
		p.Write([]byte{CharEscape, CharUnderline, 1})
	} else {
		p.Write([]byte{CharEscape, CharUnderline, 0})
	}
}
