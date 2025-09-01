package escpos

import "github.com/scott-ainsworth/go-ascii"

func (p *Printer) Emphasize(on bool) {
	if on {
		p.Write(ascii.ESC, 'E', 1)
	} else {
		p.Write(ascii.ESC, 'E', 0)
	}
}

func (p *Printer) DoubleStrike(on bool) {
	if on {
		p.Write(ascii.ESC, 'G', 1)
	} else {
		p.Write(ascii.ESC, 'G', 0)
	}
}

func (p *Printer) Underline(on bool) {
	if on {
		p.Write(ascii.ESC, '-', 1)
	} else {
		p.Write(ascii.ESC, '-', 0)
	}
}
