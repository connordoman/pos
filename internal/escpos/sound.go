package escpos

type BuzzerInstruction struct {
	Pattern   uint8
	Times     uint8
	StartTime uint8
	EndTime   uint8
}

func (p *Printer) Beep(times int) error {
	return nil
}
