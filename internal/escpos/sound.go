package escpos

import "github.com/scott-ainsworth/go-ascii"

type BuzzerInstruction struct {
	Pattern   uint8
	Times     uint8
	StartTime uint8
	EndTime   uint8
}

// Beep makes the printer emit a beep sound for `time` * 100 milliseconds
func (p *Printer) Beep(time uint8) error {
	p.Write(ascii.ESC, '(', 'A', 0x04, 0x00, 0x30, 0x00, 1, time)
	return nil
}
