package parser

import (
	"strings"

	"github.com/Clement-Jean/protein/lexer"
)

func (p *Impl) parsePackage() *string {
	if !p.acceptPeek(lexer.TokenIdentifier) {
		return nil
	}

	var parts []string

	for p.curToken.Type == lexer.TokenIdentifier {
		parts = append(parts, p.curToken.Literal)
		if p.peekToken.Type != lexer.TokenDot {
			break
		}

		p.nextToken()

		if !p.acceptPeek(lexer.TokenIdentifier) {
			return nil
		}
	}
	s := strings.Join(parts, ".")

	if !p.acceptPeek(lexer.TokenSemicolon) {
		return nil
	}

	return &s
}
