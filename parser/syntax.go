package parser

import (
	"strings"

	"github.com/Clement-Jean/protein/lexer"
)

func (p *Impl) parseSyntax() *string {
	if !p.acceptPeek(lexer.TokenEqual) {
		return nil
	}
	if !p.acceptPeek(lexer.TokenStr) {
		return nil
	}

	s := strings.TrimFunc(p.curToken.Literal, func(r rune) bool {
		return r == '\'' || r == '"'
	})

	if !p.acceptPeek(lexer.TokenSemicolon) {
		return nil
	}

	return &s
}
