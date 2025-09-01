package escpos

import (
	"errors"
	"log"
)

const (
	AFullCut     = 0
	AFullCut0    = 48
	APartialCut  = 1
	APartialCut0 = 49

	BFullCut    = 65
	BPartialCut = 66

	CFullCut    = 97
	CPartialCut = 98

	DFullCut    = 103
	DPartialCut = 104
)

func (p *Printer) SelectCutModeAndCut(function rune, m, n byte) error {
	bytes := []byte{0x1D, 0x56}
	switch function {
	case 'A':
		if m != AFullCut && m != AFullCut0 && m != APartialCut && m != APartialCut0 {
			return errors.New("invalid cut mode for function A")
		}
		bytes = append(bytes, m)
	case 'B':
		if m != BFullCut && m != BPartialCut {
			return errors.New("invalid cut mode for function B")
		}
		bytes = append(bytes, m, n)
	case 'C':
		if m != CFullCut && m != CPartialCut {
			return errors.New("invalid cut mode for function C")
		}
		bytes = append(bytes, m, n)
	case 'D':
		if m != DFullCut && m != DPartialCut {
			return errors.New("invalid cut mode for function D")
		}
		bytes = append(bytes, m, n)
	default:
		return errors.New("invalid cut function")
	}

	log.Println("Cut command:", bytes)

	p.Write(bytes)

	return nil
}

func (p *Printer) Cut() error {
	return p.SelectCutModeAndCut('A', AFullCut0, 0)
}

func (p *Printer) CutPartial() error {
	return p.SelectCutModeAndCut('A', APartialCut0, 0)
}

func (p *Printer) FeedAndCut(lines uint8) error {
	lines *= 15             // each line is approx 15 turn units
	lines = min(lines, 254) // 255 is an overflow
	return p.SelectCutModeAndCut('B', BFullCut, lines)
}

func (p *Printer) FeedAndCutPartial(lines uint8) error {
	return p.SelectCutModeAndCut('B', BPartialCut, lines)
}

func AutocutterCut(p *Printer, lines uint8) error {
	return p.SelectCutModeAndCut('C', CFullCut, lines)
}

func AutocutterCutPartial(p *Printer, lines uint8) error {
	return p.SelectCutModeAndCut('C', CPartialCut, lines)
}

func ReverseToStartAfterFeedAndCut(p *Printer, lines uint8) error {
	return p.SelectCutModeAndCut('D', DFullCut, lines)
}

func ReverseToStartAfterFeedAndCutPartial(p *Printer, lines uint8) error {
	return p.SelectCutModeAndCut('D', DPartialCut, lines)
}
