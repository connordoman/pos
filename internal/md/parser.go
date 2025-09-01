package md

import (
	"fmt"
	"log"
)

const (
	escape       = 0x1B
	bold         = 0x45
	doubleStrike = 0x47
)

func Parse(text string) ([]byte, error) {
	bytes := []byte{}

	boldCounter := 0
	doubleStrikeCounter := 0

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
				bytes = append(bytes, escape, bold, 1)
				i += 2
			} else if nextC == '_' && nextNextC == '_' {
				doubleStrikeCounter++
				bytes = append(bytes, escape, doubleStrike, 1)
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
				bytes = append(bytes, escape, bold, 0)
				i += 1
			}
		case '_':
			if doubleStrikeCounter == 0 {
				bytes = append(bytes, c)
				continue
			}

			nextC := text[nextIndex]
			nextNextC := text[nextNextIndex]
			if nextIndex == nextNextIndex {
				nextNextC = 0x00
			}
			if nextC == '_' && (nextNextC == '\n' || nextNextC == ' ' || nextNextC == 0x00) {
				doubleStrikeCounter--
				bytes = append(bytes, escape, doubleStrike, 0)
				i += 1
			}
		default:
			bytes = append(bytes, c)
		}

	}

	log.Printf("boldCounter: %d, doubleStrikeCounter: %d", boldCounter, doubleStrikeCounter)

	if boldCounter != 0 {
		return nil, fmt.Errorf("unmatched bold markers")
	}

	if doubleStrikeCounter != 0 {
		return nil, fmt.Errorf("unmatched double strike markers")
	}

	return bytes, nil
}
