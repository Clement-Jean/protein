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

func (p *impl) parseIdentifier() (ast.Identifier, error) {
	if peek := p.peek(); peek.Kind != token.KindIdentifier {
		return ast.Identifier{}, gotUnexpected(peek, token.KindIdentifier)
	}

	next := p.nextToken()
	return ast.Identifier{ID: next.ID}, nil
}

func (p *impl) parseFullyQualifiedIdentifier() (ast.Identifier, error) {
	first, err := p.parseIdentifier()

	if err != nil {
		return first, err
	}

	parts := []token.UniqueID{first.ID}

	for peek := p.peek(); peek.Kind == token.KindDot; peek = p.peek() {
		p.nextToken()
		next, err := p.parseIdentifier()

		if err != nil {
			return ast.Identifier{}, err
		}

		parts = append(parts, next.ID)
	}

	if len(parts) > 1 {
		id := p.fm.Merge(token.KindIdentifier, parts...)
		return ast.Identifier{ID: id, Parts: parts}, nil
	}
	return ast.Identifier{ID: first.ID}, nil
}
