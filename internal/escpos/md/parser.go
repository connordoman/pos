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
	for !p.IsAtEnd() {
		p.scanLine()
	}
	p.AddToken(TokenEOF, nil)
	return p.Tokens
}

func (p *Parser) IsAtEnd() bool {
	return p.current >= p.Source.Length()
}

// scanLine parses one logical line, handling headings, hrules, and inline spans.
func (p *Parser) scanLine() {
	p.start = p.current
	// Read the line up to newline or EOF while peeking for patterns
	// Detect ATX heading starting hashes at beginning of line (allow up to 3 leading spaces)
	saved := p.current
	spaces := 0
	for spaces < 3 && p.Peek() == ' ' {
		p.Advance()
		spaces++
	}
	if p.Peek() == '#' {
		level := 0
		for p.Peek() == '#' && level < 6 {
			p.Advance()
			level++
		}
		// require space after hashes unless EOL
		if p.Peek() == ' ' || p.Peek() == '\n' || p.Peek() == 0 {
			if p.Peek() == ' ' {
				p.Advance()
			}
			// Heading inline content until end of line
			p.AddToken(TokenHeadingStart, level)
			p.scanInlineUntilNewline()
			p.AddToken(TokenHeadingEnd, level)
			if p.Peek() == '\n' {
				p.Advance()
				p.AddToken(TokenNewLine, nil)
				p.line++
			}
			return
		}
	}
	// Not a heading; restore position
	p.current = saved

	// Horizontal rule: a line with 3+ dashes and optional spaces
	if p.isHRuleAhead() {
		// Consume the rest of the line
		for p.Peek() != '\n' && p.Peek() != 0 {
			p.Advance()
		}
		p.AddToken(TokenSeparator, nil)
		if p.Peek() == '\n' {
			p.Advance()
			p.AddToken(TokenNewLine, nil)
			p.line++
		}
		return
	}

	// Normal paragraph/text line: scan inline and emit newline
	p.scanInlineUntilNewline()
	if p.Peek() == '\n' {
		p.Advance()
		p.AddToken(TokenNewLine, nil)
		p.line++
	}
}

func (p *Parser) InlineCode() {

}

func (p *Parser) CodeBlock() {

}

func (p *Parser) Heading() {
	log.Println("Heading() is unused in new parser")
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
// func (p *Parser) BoldItalicTriple() {
// 	content, ok := p.parseDelimited('*', 3)
// 	if !ok {
// 		Error(p.line, "Unterminated bold+italic.")
// 		return
// 	}
// 	p.AddToken(TokenItalic, content)
// 	p.AddToken(TokenBold, content)
// }

// func (p *Parser) Bold() {
// 	content, ok := p.parseDelimited('*', 2)
// 	if !ok {
// 		Error(p.line, "Unterminated bold.")
// 		return
// 	}
// 	p.AddToken(TokenBold, content)
// }

// func (p *Parser) Italic() {
// 	content, ok := p.parseDelimited('*', 1)
// 	if !ok {
// 		Error(p.line, "Unterminated italic.")
// 		return
// 	}
// 	p.AddToken(TokenItalic, content)
// }

// func (p *Parser) Underline() {
// 	// Support nested bold inside underline: __**text**__
// 	if p.Peek() == '*' && p.Source.CharAt(p.current+1) == '*' {
// 		// consume '**'
// 		p.IncrementCursor()
// 		p.IncrementCursor()
// 		// parse bold content
// 		boldStart := p.current
// 		endBold, ok := p.findClosingRun('*', 2)
// 		if !ok {
// 			Error(p.line, "Unterminated bold inside underline.")
// 			return
// 		}
// 		content := p.SourceAt(boldStart, endBold)
// 		// move to end of bold
// 		for p.current < endBold {
// 			p.Advance()
// 		}
// 		// advance past '**'
// 		p.IncrementCursor()
// 		p.IncrementCursor()
// 		// require closing '__'
// 		if p.Source.CharAt(p.current) == '_' && p.Source.CharAt(p.current+1) == '_' {
// 			p.IncrementCursor()
// 			p.IncrementCursor()
// 			// Emit both tokens with same content
// 			// p.AddToken(TokenBold, content)
// 			// p.AddToken(TokenUnderline, content)
// 			return
// 		}
// 		Error(p.line, "Unterminated underline.")
// 		return
// 	}

// 	// Plain underline: __text__
// 	content, ok := p.parseDelimited('_', 2)
// 	if !ok {
// 		Error(p.line, "Unterminated underline.")
// 		return
// 	}
// 	p.AddToken(TokenUnderline, content)
// }

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
	// Lexeme here is best-effort; for inline text we’ll override when adding text tokens
	var lex string
	if p.current > p.start {
		lex = p.Source[p.start:p.current].String()
	}
	p.Tokens = append(p.Tokens, NewToken(t, lex, obj, p.line))
}

// Helpers for the new parser

func (p *Parser) isHRuleAhead() bool {
	// Allow leading spaces
	i := p.current
	for p.Source.CharAt(i) == ' ' {
		i++
	}
	// Count dashes
	dashes := 0
	for p.Source.CharAt(i) == '-' {
		dashes++
		i++
	}
	if dashes < 3 {
		return false
	}
	// Only spaces until newline or end
	for {
		r := p.Source.CharAt(i)
		if r == 0 || r == '\n' {
			return true
		}
		if r != ' ' && r != '-' { // Allow extra dashes/spaces
			return false
		}
		i++
	}
}

func (p *Parser) scanInlineUntilNewline() {
	// Emit text and formatting events until newline or EOF
	var buf []rune
	flushText := func() {
		if len(buf) > 0 {
			s := string(buf)
			p.Tokens = append(p.Tokens, NewToken(TokenText, s, s, p.line))
			buf = buf[:0]
		}
	}
	for {
		r := p.Peek()
		if r == 0 || r == '\n' {
			flushText()
			return
		}
		if r == '\\' {
			// Escape next char literally
			p.Advance()
			next := p.Peek()
			if next == 0 || next == '\n' {
				buf = append(buf, '\\')
				continue
			}
			p.Advance()
			buf = append(buf, next)
			continue
		}
		// Bold with **
		if r == '*' && p.Source.CharAt(p.current+1) == '*' {
			flushText()
			// toggle bold
			p.Advance() // first *
			p.Advance() // second *
			// Peek ahead: if we are at a closing (next non-space is not delimiter start?), we cannot easily know; we’ll parse balanced by scanning content until next ** when emitting at consumer. Simpler: emit a BoldStart token and rely on consumer to stack. But we need paired events. Here we choose to treat encountering ** as a toggle; consumer maintains a stack.
			// To enable pairing, we will check if we are currently inside bold by looking at the last token; however parser shouldn’t track state. We'll instead look ahead for a matching ** before newline; if found, emit start, parse text up to closing, then emit end.
			// Find closing **
			if idx, ok := p.findClosingRun('*', 2); ok {
				// Emit start
				p.Tokens = append(p.Tokens, NewToken(TokenBoldStart, "**", nil, p.line))
				// Emit inner as text with further inline handling (allow nested underline). We'll recursively scan by slicing temporarily — but to keep it simple, parse inline within the range manually.
				// We'll consume until idx, but process underscores escapes inside.
				for p.current < idx {
					rr := p.Peek()
					if rr == '\\' {
						p.Advance()
						n2 := p.Peek()
						if n2 == 0 || n2 == '\n' {
							buf = append(buf, '\\')
							continue
						}
						p.Advance()
						buf = append(buf, n2)
						continue
					}
					if rr == '_' { // underline single underscore
						flushText()
						p.Advance()
						// find closing _
						if cidx, ok2 := p.findClosingRun('_', 1); ok2 && cidx <= idx {
							p.Tokens = append(p.Tokens, NewToken(TokenUnderlineStart, "_", nil, p.line))
							// collect text until cidx
							for p.current < cidx {
								chr := p.Advance()
								buf = append(buf, chr)
							}
							flushText()
							// consume closing _
							p.Advance()
							p.Tokens = append(p.Tokens, NewToken(TokenUnderlineEnd, "_", nil, p.line))
							continue
						}
						// no closing, treat as literal
						buf = append(buf, '_')
						continue
					}
					// regular char inside bold segment
					buf = append(buf, p.Advance())
				}
				flushText()
				// consume closing **
				p.Advance()
				p.Advance()
				p.Tokens = append(p.Tokens, NewToken(TokenBoldEnd, "**", nil, p.line))
				continue
			}
			// No closing **; treat as literal '**'
			buf = append(buf, '*', '*')
			continue
		}
		if r == '_' { // single underscore underline
			flushText()
			p.Advance()
			if idx, ok := p.findClosingRun('_', 1); ok {
				p.Tokens = append(p.Tokens, NewToken(TokenUnderlineStart, "_", nil, p.line))
				for p.current < idx {
					buf = append(buf, p.Advance())
				}
				flushText()
				p.Advance() // consume closing _
				p.Tokens = append(p.Tokens, NewToken(TokenUnderlineEnd, "_", nil, p.line))
				continue
			}
			// no closing; literal
			buf = append(buf, '_')
			continue
		}
		// Regular character: accumulate
		buf = append(buf, p.Advance())
	}
}
