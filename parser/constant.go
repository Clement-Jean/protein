package parser

import (
	"github.com/Clement-Jean/protein/ast"
	"github.com/Clement-Jean/protein/token"
)

func (p *impl) parseString() (ast.String, error) {
	if peek := p.peek(); peek.Kind != token.KindStr {
		return ast.String{}, gotUnexpected(peek, token.KindStr)
	}

	next := p.nextToken()
	return ast.String{ID: next.ID}, nil
}
