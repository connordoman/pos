package escpos

import "fmt"

const (
	ContinuousLineThin  = 1
	ContinuousLineThick = 2
	ContinuousLineBlack = 3
)

// Aliases matching the ESC/POS spec names for m1 (line type/thickness)
const (
	LineThin   = ContinuousLineThin  // m1 = 1
	LineMedium = ContinuousLineThick // m1 = 2 (moderately thick)
	LineThick  = ContinuousLineBlack // m1 = 3 (thick)
)

// Color (c) parameter. Most TM‑T88 series (and many Rongta clones) are monochrome;
// use ColorPrimary. If your device supports 2‑color thermal paper, ColorSecondary may draw in the second color.
const (
	ColorPrimary   byte = 0x00
	ColorSecondary byte = 0x01
)

// Mode (m2) parameter. For TM‑T88V / TM‑T88VII class devices this is typically 0x00.
const DefaultLineMode byte = 0x00

// DrawLinePage enqueues GS ( Q <fn=48> to draw a horizontal or vertical line in Page mode.
// Coordinates are in printer motion units (typically dots). Only horizontal (y1==y2)
// or vertical (x1==x2) lines are valid.
//
// c  = color selector (use ColorPrimary for monochrome devices)
// m1 = thickness (1: thin, 2: medium, 3: thick)
// m2 = extra mode (usually 0)
func (p *Printer) DrawLinePage(x1, y1, x2, y2 uint16, c, m1, m2 byte) error {
	// Validate constraints that would otherwise be ignored by the device.
	if x1 == x2 && y1 == y2 {
		return fmt.Errorf("draw line: start and end coordinates are identical")
	}
	if !(x1 == x2 || y1 == y2) {
		return fmt.Errorf("draw line: only horizontal or vertical lines are supported")
	}

	// Command header: GS ( Q, length=12 (0x0C,0x00), fn=0x30 (48)
	bytes := []byte{0x1D, 0x28, 0x51, 0x0C, 0x00, 0x30}

	// x1, y1, x2, y2 in little‑endian (L then H)
	bytes = append(bytes,
		byte(x1&0xFF), byte(x1>>8),
		byte(y1&0xFF), byte(y1>>8),
		byte(x2&0xFF), byte(x2>>8),
		byte(y2&0xFF), byte(y2>>8),
		c, m1, m2,
	)

	p.Write(bytes)
	return nil
}

// DrawLineSimple uses safe defaults for TM‑T88V/TM‑T88VII‑class printers (and many Rongta clones):
// c=ColorPrimary, m2=0.
func (p *Printer) DrawLineSimple(x1, y1, x2, y2 uint16, thickness byte) error {
	return p.DrawLinePage(x1, y1, x2, y2, ColorPrimary, thickness, DefaultLineMode)
}
