package escpos

import (
	"fmt"

	"github.com/connordoman/pos/internal/escpos/md"
)

type ASTNode struct {
	Type     md.TokenType
	Literal  string
	Children []*ASTNode
}

func NewAST() *ASTNode {
	root := &ASTNode{
		Type:     md.TokenEOF,
		Literal:  "",
		Children: []*ASTNode{},
	}
	return root
}

func (n *ASTNode) AddChild(child *ASTNode) {
	n.Children = append(n.Children, child)
}

func (p *Printer) ParseMarkdown(text string) error {
	interpreter := md.NewInterpreter()
	err := interpreter.Run(text)
	if err != nil {
		return fmt.Errorf("failed to parse markdown: %w", err)
	}

	// Walk tokens and emit ESC/POS
	bold := false
	underline := false
	inHeading := false
	for _, t := range interpreter.Tokens {
		switch t.Type {
		case md.TokenHeadingStart:
			// For headings, just enable emphasize; ignore level for now
			p.Emphasize(true)
			inHeading = true
		case md.TokenHeadingEnd:
			if inHeading {
				p.Emphasize(false)
				inHeading = false
			}
		case md.TokenBoldStart:
			if !bold {
				p.Emphasize(true)
				bold = true
			}
		case md.TokenBoldEnd:
			if bold {
				p.Emphasize(false)
				bold = false
			}
		case md.TokenUnderlineStart:
			if !underline {
				p.Underline(true)
				underline = true
			}
		case md.TokenUnderlineEnd:
			if underline {
				p.Underline(false)
				underline = false
			}
		case md.TokenText:
			if s, ok := t.Literal.(string); ok {
				p.WriteString(s)
			}
		case md.TokenNewLine:
			p.WriteString("\n")
		case md.TokenSeparator:
			p.SimpleLine()
		case md.TokenCodeBlock:
			// Treat as plain text for now (already included as TokenText in new parser), ignore
		case md.TokenEOF:
			// no-op
		default:
			// ignore unknowns safely
		}
	}

	return nil
}
