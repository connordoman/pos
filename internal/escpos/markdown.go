package escpos

import (
	"fmt"
	"strings"

	"github.com/connordoman/pos/internal/escpos/md"
)

func (p *Printer) ParseMarkdown(text string) error {
	interpreter := md.NewInterpreter()
	err := interpreter.Run(text)
	if err != nil {
		return fmt.Errorf("failed to parse markdown: %w", err)
	}

	for _, token := range interpreter.Tokens {
		switch token.Type {
		case md.TokenHeading1, md.TokenHeading2, md.TokenHeading3, md.TokenHeading4, md.TokenHeading5, md.TokenHeading6:
			mark := strings.Repeat("#", token.Type.HeadingSize())
			p.Emphasize(true)
			p.WriteString(fmt.Sprintf("%s %s", mark, token.Lexeme))
			p.Emphasize(false)
		case md.TokenBold:
			p.Emphasize(true)
			p.WriteString(token.Lexeme)
			p.Emphasize(false)
		case md.TokenUnderline:
			p.Underline(true)
			p.WriteString(token.Lexeme)
			p.Underline(false)

		}
	}

	return nil
}
