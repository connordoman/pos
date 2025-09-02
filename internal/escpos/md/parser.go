package md

import (
	"log"
)

type Parser struct {
	Source  Source
	Tokens  []*Token
	start   int
	current int
	line    int

	isNewLine bool
}

func NewParser(source string) *Parser {
	return &Parser{
		Source: Source(source),
		Tokens: []*Token{},

		isNewLine: true,
	}
}

func (p *Parser) ScanTokens(source string) []*Token {
	for {
		if p.IsAtEnd() {
			log.Println("End of input")
			break
		}

		p.start = p.current
		p.ScanToken()
	}

	p.Tokens = append(p.Tokens, NewToken(TokenEOF, "", nil, p.line))
	return p.Tokens
}

func (p *Parser) IsAtEnd() bool {
	return p.current >= p.Source.Length()
}

func (p *Parser) ScanToken() {
	c := p.Advance()

	// log.Printf("scanning token: %q", c)

	newLine := p.line

	if p.isNewLine {
		for p.Peek() == ' ' && !p.IsAtEnd() {
			p.Advance()
		}

		// Error(p.line, "Unexpected character.")
	}

	switch c {
	case '#':
		if p.isNewLine {
			p.Heading()
			newLine++
		}
	// case '`':
	// 	p.InlineCode()
	case '*':
		// If escaped, treat as literal and skip
		if p.Prev() == '\\' {
			break
		}
		// Determine if this is ** (bold) or *** (bold+italic) or * (italic)
		if p.Match('*') {
			// We consumed a second '*', check for triple '***'
			if p.Match('*') {
				p.BoldItalicTriple()
			} else {
				p.Bold()
			}
		} else {
			p.Italic()
		}
	case '_':
		// If escaped, treat as literal and skip
		if p.Prev() == '\\' {
			break
		}
		if p.Match('_') {
			p.Underline()
		}
	case '\n':
		newLine++
	default:
	}

	if newLine > p.line {
		p.isNewLine = true
		p.line = newLine
	} else if p.isNewLine {
		p.isNewLine = false
	}

	// log.Printf("c: %q (isNewLine=%t)", c, p.isNewLine)
}

func (p *Parser) InlineCode() {

}

func (p *Parser) CodeBlock() {

}

func (p *Parser) Heading() {
	log.Println("Found heading")
	headingSize := 1

	for p.Match('#') {
		if headingSize < 6 {
			headingSize++
		}
	}

	for p.Match(' ') {
	}
	p.Advance()

	p.start = p.current - 1

	var headingContent string = p.ConsumeLine()

	switch headingSize {
	case 1:
		p.AddToken(TokenHeading1, headingContent)
	case 2:
		p.AddToken(TokenHeading2, headingContent)
	case 3:
		p.AddToken(TokenHeading3, headingContent)
	case 4:
		p.AddToken(TokenHeading4, headingContent)
	case 5:
		p.AddToken(TokenHeading5, headingContent)
	case 6:
		p.AddToken(TokenHeading6, headingContent)
	}
}

// Prev returns the previous rune (or NUL if not available)
func (p *Parser) Prev() rune {
	if p.current-2 < 0 {
		return 0
	}
	return p.Source.CharAt(p.current - 2)
}

// runLenAt counts consecutive occurrences of ch starting at absolute index i
func (p *Parser) runLenAt(i int, ch rune) int {
	n := 0
	for {
		r := p.Source.CharAt(i + n)
		if r == ch && r != 0 {
			n++
			continue
		}
		break
	}
	return n
}

// findClosingRun finds the index of the start of a closing run of `count` delimiters `ch`,
// honoring simple escapes (\\). Returns the index of the first closing delimiter and ok.
func (p *Parser) findClosingRun(ch rune, count int) (int, bool) {
	i := p.current
	for !p.IsAtEnd() && p.Source.CharAt(i) != 0 {
		r := p.Source.CharAt(i)
		if r == '\\' {
			// Skip escaped next character
			i += 2
			continue
		}
		if p.runLenAt(i, ch) >= count {
			return i, true
		}
		i++
	}
	return 0, false
}

// parseDelimited assumes the opening delimiter (of given count) has already been consumed.
// It captures content up to a matching closing delimiter sequence and advances the cursor past it.
func (p *Parser) parseDelimited(ch rune, count int) (string, bool) {
	start := p.current
	end, ok := p.findClosingRun(ch, count)
	if !ok {
		return "", false
	}
	content := p.SourceAt(start, end)
	// Move cursor to end
	for p.current < end {
		p.Advance()
	}
	// Advance current past the closing run
	for i := 0; i < count; i++ {
		p.Advance()
	}
	return content, true
}

// BoldItalicTriple handles ***text*** by emitting Italic then Bold tokens over the same content
func (p *Parser) BoldItalicTriple() {
	content, ok := p.parseDelimited('*', 3)
	if !ok {
		Error(p.line, "Unterminated bold+italic.")
		return
	}
	p.AddToken(TokenItalic, content)
	p.AddToken(TokenBold, content)
}

func (p *Parser) Bold() {
	content, ok := p.parseDelimited('*', 2)
	if !ok {
		Error(p.line, "Unterminated bold.")
		return
	}
	p.AddToken(TokenBold, content)
}

func (p *Parser) Italic() {
	content, ok := p.parseDelimited('*', 1)
	if !ok {
		Error(p.line, "Unterminated italic.")
		return
	}
	p.AddToken(TokenItalic, content)
}

func (p *Parser) Underline() {
	// Support nested bold inside underline: __**text**__
	if p.Peek() == '*' && p.Source.CharAt(p.current+1) == '*' {
		// consume '**'
		p.IncrementCursor()
		p.IncrementCursor()
		// parse bold content
		boldStart := p.current
		endBold, ok := p.findClosingRun('*', 2)
		if !ok {
			Error(p.line, "Unterminated bold inside underline.")
			return
		}
		content := p.SourceAt(boldStart, endBold)
		// move to end of bold
		for p.current < endBold {
			p.Advance()
		}
		// advance past '**'
		p.Advance()
		p.Advance()
		// require closing '__'
		if p.Source.CharAt(p.current) == '_' && p.Source.CharAt(p.current+1) == '_' {
			p.IncrementCursor()
			p.IncrementCursor()
			// Emit both tokens with same content
			p.AddToken(TokenBold, content)
			p.AddToken(TokenUnderline, content)
			return
		}
		Error(p.line, "Unterminated underline.")
		return
	}

	// Plain underline: __text__
	content, ok := p.parseDelimited('_', 2)
	if !ok {
		Error(p.line, "Unterminated underline.")
		return
	}
	p.AddToken(TokenUnderline, content)
}

func (p *Parser) ConsumeLine() string {
	// p.Advance()
	for !p.Match('\n') && !p.IsAtEnd() {
		p.Advance()
	}
	return p.CurrentLexeme()
}

func (p *Parser) SourceAt(start, end int) string {
	return p.Source.Substring(start, end)
}

func (p *Parser) CurrentLexeme() string {
	return p.SourceAt(p.start, p.current-1)
}

// func (p *Parser) text() {
// 	for p.peek() != '\n' && !p.IsAtEnd() {
// 		p.advance()
// 	}

// 	if p.IsAtEnd() {
// 		// Error(p.line, "Unterminated string.")
// 	}

// 	value := p.Source.Substring(p.start, p.current)
// 	p.addToken(TokenString, value)
// }

func (p *Parser) Match(expected rune) bool {
	if p.IsAtEnd() {
		return false
	}
	if p.Peek() != expected {
		return false
	}

	p.IncrementCursor()
	return true
}

func (p *Parser) Peek() rune {
	if p.IsAtEnd() {
		return 0x00
	}
	return p.Source.CharAt(p.current)
}

func (p *Parser) Advance() rune {
	char := p.Source.CharAt(p.current)
	p.IncrementCursor()
	return char
}

func (p *Parser) IncrementCursor() {
	p.current++
}

func (p *Parser) AddNullToken(t TokenType) {
	p.AddToken(t, nil)
}

func (p *Parser) AddToken(t TokenType, obj any) {
	text := p.Source[p.start:p.current]
	p.Tokens = append(p.Tokens, NewToken(t, text.String(), obj, p.line))
}
