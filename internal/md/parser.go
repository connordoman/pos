package md

import (
	"fmt"
	"log"

	"github.com/connordoman/pos/internal/escpos"
)

func Parse(text string) ([]byte, error) {
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
				bytes = append(bytes,
					escpos.CharEscape,
					escpos.CharBold,
					1,
				)
				i += 2
			} else if nextC == '_' && nextNextC == '_' {
				underlineCounter++
				bytes = append(bytes,
					escpos.CharEscape,
					escpos.CharUnderline,
					1,
				)
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
				bytes = append(bytes,
					escpos.CharEscape,
					escpos.CharBold,
					0,
				)
				i += 1
			}
		case '_':
			if underlineCounter == 0 {
				bytes = append(bytes, c)
				continue
			}

			nextC := text[nextIndex]
			nextNextC := text[nextNextIndex]
			if nextIndex == nextNextIndex {
				nextNextC = 0x00
			}
			if nextC == '_' && (nextNextC == '\n' || nextNextC == ' ' || nextNextC == 0x00) {
				underlineCounter--
				bytes = append(bytes,
					escpos.CharEscape,
					escpos.CharUnderline,
					0,
				)
				i += 1
			}
		default:
			bytes = append(bytes, c)
		}

	}

	log.Printf("boldCounter: %d, underlineCounter: %d", boldCounter, underlineCounter)

	if boldCounter != 0 {
		return nil, fmt.Errorf("unmatched bold markers")
	}

	if underlineCounter != 0 {
		return nil, fmt.Errorf("unmatched underline markers")
	}

	return bytes, nil
}
