package md

import "fmt"

type TokenType int

const (
	TokenEOF TokenType = iota

	TokenInlineCode
	TokenCodeBlock

	TokenBold
	TokenUnderscored

	TokenHeading1
	TokenHeading2
	TokenHeading3
	TokenHeading4
	TokenHeading5
	TokenHeading6

	TokenSeparator
)

func (t TokenType) String() string {
	switch t {
	case TokenEOF:
		return "EOF"
	case TokenInlineCode:
		return "InlineCode"
	case TokenCodeBlock:
		return "CodeBlock"
	case TokenBold:
		return "Bold"
	case TokenUnderscored:
		return "Underscored"
	case TokenHeading1:
		return "H1"
	case TokenHeading2:
		return "H2"
	case TokenHeading3:
		return "H3"
	case TokenHeading4:
		return "H4"
	case TokenHeading5:
		return "H5"
	case TokenHeading6:
		return "H6"
	case TokenSeparator:
		return "Separator"
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
		return fmt.Sprintf("%d %s %q", t.Line, t.Type, t.Literal)
	}
	return fmt.Sprintf("%d %s %v", t.Line, t.Type, t.Literal)
}
