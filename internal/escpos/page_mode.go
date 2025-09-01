package escpos

// Page mode helpers per ESC/POS spec.

// EnterPageMode sends ESC L to switch to Page mode.
func (p *Printer) EnterPageMode() {
	p.Write([]byte{0x1B, 0x4C})
}

// StandardMode sends ESC S to switch back to Standard mode.
func (p *Printer) StandardMode() {
	p.Write([]byte{0x1B, 0x53})
}

// SetPageArea sets the page print area in Page mode: ESC W xL xH yL yH dxL dxH dyL dyH
// x,y are the origin; dx,dy are the width/height in motion units (typically dots).
func (p *Printer) SetPageArea(x, y, dx, dy uint16) {
	bytes := []byte{0x1B, 0x57,
		byte(x & 0xFF), byte(x >> 8),
		byte(y & 0xFF), byte(y >> 8),
		byte(dx & 0xFF), byte(dx >> 8),
		byte(dy & 0xFF), byte(dy >> 8),
	}
	p.Write(bytes)
}

// PrintPage prints the contents of the page buffer: ESC FF.
func (p *Printer) PrintPage() {
	p.Write([]byte{0x1B, 0x0C})
}
