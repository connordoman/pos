package escpos

func (p *Printer) Emphasize(on bool) {
	if on {
		p.Write([]byte{0x1B, 0x45, 1})
	} else {
		p.Write([]byte{0x1B, 0x45, 0})
	}
}

func (p *Printer) DoubleStrike(on bool) {
	if on {
		p.Write([]byte{0x1B, 0x47, 1})
	} else {
		p.Write([]byte{0x1B, 0x47, 0})
	}
}
