package escpos

import (
	"fmt"
	"log"
)

func (p *Printer) ParseMarkdown(text string) error {
	bytes := []byte{}

	boldCounter := 0
	underlineCounter := 0

	for i := 0; i < len(text); i++ {
		c := text[i]

		nextIndex := min(i+1, len(text)-1)
		nextNextIndex := min(i+2, len(text)-1)

		fmt.Printf("%c", c)

		switch c {
		case ' ':
			bytes = append(bytes, c)

			nextC := text[nextIndex]
			nextNextC := text[nextNextIndex]
			if nextC == '*' && nextNextC == '*' {
				boldCounter++
				p.Emphasize(true)
				i += 2
			} else if nextC == '_' && nextNextC == '_' {
				underlineCounter++
				p.Underline(true)
				i += 2
			}

		case '*':
			if boldCounter == 0 {
				bytes = append(bytes, c)
				continue
			}

			nextC := text[nextIndex]
			nextNextC := text[nextNextIndex]
			if nextIndex == nextNextIndex {
				nextNextC = 0x00
			}
			if nextC == '*' && (nextNextC == '\n' || nextNextC == ' ' || nextNextC == 0x00) {
				boldCounter--
				p.Emphasize(false)
				i += 1
			}
		case '_':
			if underlineCounter == 0 {
				p.Write(c)
				continue
			}

			nextC := text[nextIndex]
			nextNextC := text[nextNextIndex]
			if nextIndex == nextNextIndex {
				nextNextC = 0x00
			}
			if nextC == '_' && (nextNextC == '\n' || nextNextC == ' ' || nextNextC == 0x00) {
				underlineCounter--
				p.Underline(false)
				i += 1
			}
		default:
			bytes = append(bytes, c)
		}

	}

	log.Printf("boldCounter: %d, underlineCounter: %d", boldCounter, underlineCounter)

	if boldCounter != 0 {
		return fmt.Errorf("unmatched bold markers")
	}

	if underlineCounter != 0 {
		return fmt.Errorf("unmatched underline markers")
	}

	p.Write(bytes...)

	return nil
}
