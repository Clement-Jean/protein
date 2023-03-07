package parser

import (
	"github.com/Clement-Jean/protein/lexer"
)

func (p *Impl) parseImport() string {
	if !p.acceptPeek(lexer.TokenStr) {
		return ""
	}

	s := destringify(p.curToken.Literal)

	if !p.acceptPeek(lexer.TokenSemicolon) {
		return ""
	}

	return s
}
