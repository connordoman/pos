package escpos

import "strings"

func substituteUnicode(s string) string {
	var b strings.Builder

	for _, r := range s {
		c := 0x00

		switch r {
		case '\u2018', '\u2019': // single quotation marks
			c = '\''
		case '\u201c', '\u201d': // double quotation marks
			c = '"'
		case '\u2013', '\u2014': // en dash and em dash
			c = '-'
		}

		if c == 0x00 && r > 0x7F {
			continue
		}

		if c != 0x00 {
			b.WriteRune(rune(c))
		} else {
			b.WriteRune(r)
		}
	}

	return b.String()
}
