package md

import "log"

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
		for p.Match(' ') && !p.IsAtEnd() {
			// p.Advance()
		}

		// Error(p.line, "Unexpected character.")
	}

	switch c {
	// single tokens
	case '#':
		log.Println("Found #")
		if p.isNewLine {
			log.Println("At start of line, parsing heading")
			p.Heading()
		}
	case '\n':
		newLine++
	default:
		p.Advance()
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
	return p.SourceAt(p.start, p.current)
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
