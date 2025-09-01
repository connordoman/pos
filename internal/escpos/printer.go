package escpos

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/gousb"
)

const (
	VendorId  = 0x0fe6
	ProductId = 0x811e

	CharactersPerLineFontA = 48
	CharactersPerLineFontB = 64
	CharactersPerLineFontC = 64

	FontAName = 'A'
	FontBName = 'B'
	FontCName = 'C'

	FontA  = 0
	FontA0 = 48

	FontB  = 1
	FontB0 = 49

	FontC  = 2
	FontC0 = 50

	SpecialFontA = 97
	SpecialFontB = 98

	UseNought = true
)

type Printer struct {
	// Keep references alive for the duration of the printer usage.
	ctx    *gousb.Context
	device *gousb.Device
	cfg    *gousb.Config
	intf   *gousb.Interface
	out    *gousb.OutEndpoint
	buff   []byte

	fontName rune

	useNought bool
}

func InitPrinter() (*Printer, error) {
	p := &Printer{
		ctx:       nil,
		device:    nil,
		cfg:       nil,
		intf:      nil,
		out:       nil,
		buff:      make([]byte, 0, 256),
		fontName:  'A',
		useNought: false,
	}

	ctx := gousb.NewContext()

	dev, err := ctx.OpenDeviceWithVIDPID(VendorId, ProductId)
	if err != nil {
		ctx.Close()
		return nil, err
	}
	if dev == nil {
		ctx.Close()
		return nil, errors.New("device not found")
	}
	p.ctx = ctx
	p.device = dev

	dev.SetAutoDetach(true)

	cfg, err := p.device.Config(1)
	if err != nil {
		p.Close()
		return nil, err
	}
	p.cfg = cfg

	intf, err := cfg.Interface(0, 0)
	if err != nil {
		p.Close()
		return nil, err
	}
	p.intf = intf

	ep, err := intf.OutEndpoint(1)
	if err != nil {
		p.Close()
		return nil, err
	}

	p.out = ep

	return p, nil
}

func (p *Printer) Init() error {
	init := []byte{0x1b, 0x40}
	p.Write(init)
	return nil
}

func (p *Printer) SelectFont(name rune) error {
	b := FontA
	switch name {
	case FontAName:
		if p.useNought {
			b = FontA0
		} else {
			b = FontA
		}
	case FontBName:
		if p.useNought {
			b = FontB0
		} else {
			b = FontB
		}
	case FontCName:
		if p.useNought {
			b = FontC0
		} else {
			b = FontC
		}
	default:
		return fmt.Errorf("invalid font name: %c", name)
	}
	nameByte := []byte{0x1B, 0x4D, byte(b)}
	p.Write(nameByte)
	p.fontName = name
	return nil
}

func (p *Printer) Log(message string) {
	t := time.Now().Format("2006-01-02 15:04:05 MST")
	s := fmt.Sprintf("%s %s", t, message)
	log.Println(s)
	p.WriteString(s + "\n")
}

func (p Printer) CharactersPerLine() int {
	switch p.fontName {
	case FontAName:
		return CharactersPerLineFontA
	case FontBName:
		return CharactersPerLineFontB
	case FontCName:
		return CharactersPerLineFontC
	default:
		return 0
	}
}

func (p *Printer) SimpleLine() {
	var builder strings.Builder
	for range p.CharactersPerLine() {
		builder.WriteByte(0xC4)
	}
	line := builder.String()
	p.WriteString(line + "\n")
}

func (p *Printer) Write(data []byte) {
	p.buff = append(p.buff, data...)
}

func (p *Printer) Flush() (int, error) {
	if p.out == nil {
		return 0, errors.New("output endpoint not initialized")
	}
	if len(p.buff) == 0 {
		return 0, nil
	}
	// Some printers behave differently when many commands are coalesced.
	// Write in reasonably sized chunks to mimic multiple writes.
	const chunkSize = 512
	written := 0
	for off := 0; off < len(p.buff); off += chunkSize {
		end := off + chunkSize

		end = min(end, len(p.buff))

		output := p.buff[off:end]
		log.Println("Flushing bytes:", output)
		n, err := p.out.Write(output)
		written += n
		if err != nil {
			// Keep remaining bytes; caller may retry or inspect.
			p.buff = p.buff[off+n:]
			return written, err
		}
		if n != end-off {
			// Short write without error: treat as error condition.
			p.buff = p.buff[off+n:]
			return written, errors.New("short write to USB endpoint")
		}
	}
	// All sent, clear buffer
	p.buff = p.buff[:0]
	return written, nil
}

func (p *Printer) Close() error {
	// Drop buffer reference
	p.buff = nil
	if p.intf != nil {
		p.intf.Close()
		p.intf = nil
	}
	if p.cfg != nil {
		p.cfg.Close()
		p.cfg = nil
	}
	if p.device != nil {
		p.device.Close()
		p.device = nil
	}
	if p.ctx != nil {
		p.ctx.Close()
		p.ctx = nil
	}
	return nil
}

func (p *Printer) WriteString(s string) {
	p.Write([]byte(s))
}

func (p *Printer) Feed(lines uint8) {
	p.Write([]byte{0x1B, 0x64, lines})
}

func (p *Printer) TestPrint() {
	// https://download4.epson.biz/sec_pubs/pos/reference_en/escpos/gs_lparen_ca.html
	// paper, n = [(0, 48) = "Basic sheet (roll paper)"],
	// 			  [(1, 49), (2, 50) = "Roll paper"],
	n := byte(0x01)
	// test pattern, m = [(1, 49) = "Hex dump"],
	// 					 [(2, 50) = "Print status printing"],
	// 					 [(3, 51) = "Rolling pattern"],
	// 					 [64 = "Automatic setting of paper layout"],
	m := byte(0x01)
	test := []byte{
		0x1D, 0x28, 0x41, 0x02, 0x00, n, m,
	}
	p.Write(test)
}
