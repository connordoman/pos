package md

import "fmt"

type TokenType int

const (
	TokenEOF TokenType = iota

	// Structural
	TokenNewLine
	TokenSeparator // horizontal rule

	// Headings (line-level). We’ll emit a HeadingStart(level) and HeadingEnd to allow inline formatting inside
	TokenHeadingStart // Literal: int level (1..6)
	TokenHeadingEnd

	// Inline formatting events (hierarchical start/end toggles)
	TokenBoldStart
	TokenBoldEnd
	TokenUnderlineStart
	TokenUnderlineEnd

	// Inline text
	TokenText

	// Code constructs (optional; treat as plain text for now if encountered)
	TokenCodeBlock
)

func (t TokenType) String() string {
	switch t {
	case TokenEOF:
		return "EOF"
	case TokenNewLine:
		return "NewLine"
	case TokenSeparator:
		return "Separator"
	case TokenHeadingStart:
		return "HeadingStart"
	case TokenHeadingEnd:
		return "HeadingEnd"
	case TokenBoldStart:
		return "BoldStart"
	case TokenBoldEnd:
		return "BoldEnd"
	case TokenUnderlineStart:
		return "UnderlineStart"
	case TokenUnderlineEnd:
		return "UnderlineEnd"
	case TokenText:
		return "Text"
	case TokenCodeBlock:
		return "CodeBlock"
	default:
		return "Unknown"
	}
}

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal any
	Line    int
}

func NewToken(t TokenType, lexeme string, literal any, line int) *Token {
	return &Token{
		Type:    t,
		Lexeme:  lexeme,
		Literal: literal,
		Line:    line,
	}
}

func (t *Token) String() string {
	switch t.Literal.(type) {
	case string:
		return fmt.Sprintf("%4d %10s %q", t.Line, t.Type, t.Literal)
	}
	return fmt.Sprintf("%4d %10s %v", t.Line, t.Type, t.Literal)
}
