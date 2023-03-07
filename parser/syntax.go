package parser

import (
	"github.com/Clement-Jean/protein/lexer"
)

func (p *Impl) parseSyntax() *string {
	if !p.acceptPeek(lexer.TokenEqual) {
		return nil
	}
	if !p.acceptPeek(lexer.TokenStr) {
		return nil
	}

	s := destringify(p.curToken.Literal)

	if !p.acceptPeek(lexer.TokenSemicolon) {
		return nil
	}

	return &s
}
