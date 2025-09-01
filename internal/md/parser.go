package md

import "fmt"

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

		switch c {
		case ' ':
			bytes = append(bytes, c)

			nextC := text[min(i+1, len(text)-1)]
			nextNextC := text[min(i+2, len(text)-1)]
			if nextC == '*' && nextNextC == '*' {
				boldCounter++
				bytes = append(bytes, escape, bold, 1)
				i += 1
			} else if nextC == '_' && nextNextC == '_' {
				doubleStrikeCounter++
				bytes = append(bytes, escape, doubleStrike, 1)
				i += 1
			}

		case '*':
			nextC := text[min(i+1, len(text)-1)]
			nextNextC := text[min(i+2, len(text)-1)]
			if nextC == '*' && (nextNextC == ' ' || nextNextC == 0x00) {
				boldCounter--
				bytes = append(bytes, escape, bold, 0)
				i += 1
			}
		case '_':
			nextC := text[min(i+1, len(text)-1)]
			nextNextC := text[min(i+2, len(text)-1)]
			if nextC == '_' && (nextNextC == ' ' || nextNextC == 0x00) {
				doubleStrikeCounter--
				bytes = append(bytes, escape, doubleStrike, 0)
				i += 1
			}
		}

	}

	if boldCounter != 0 {
		return nil, fmt.Errorf("unmatched bold markers")
	}

	if doubleStrikeCounter != 0 {
		return nil, fmt.Errorf("unmatched double strike markers")
	}

	return bytes, nil
}
