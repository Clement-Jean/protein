package parser

import (
	"github.com/Clement-Jean/protein/ast"
	"github.com/Clement-Jean/protein/token"
)

func (p *impl) parsePackage() (ast.Package, error) {
	first := p.curr()
	lit, err := p.parseFullyQualifiedIdentifier()

	if err != nil {
		return ast.Package{}, err
	}

	if peek := p.peek(); peek.Kind != token.KindSemicolon {
		return ast.Package{}, gotUnexpected(peek, token.KindSemicolon)
	}

	last := p.nextToken()
	id := p.fm.Merge(token.KindPackage, first.ID, last.ID)
	return ast.Package{ID: id, Value: lit}, nil
}
