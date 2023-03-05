package parser

import (
	"strings"

	"github.com/Clement-Jean/protein/lexer"
)

func (p *Impl) parseSyntax() *string {
	if !p.expectPeek(lexer.TokenEqual) {
		return nil
	}
	if !p.expectPeek(lexer.TokenStr) {
		return nil
	}

	s := strings.TrimFunc(p.curToken.Literal, func(r rune) bool {
		return r == '\'' || r == '"'
	})

	if !p.expectPeek(lexer.TokenSemicolon) {
		return nil
	}

	return &s
}
