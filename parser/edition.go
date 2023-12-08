package parser

import (
	"github.com/Clement-Jean/protein/ast"
	"github.com/Clement-Jean/protein/token"
)

func (p *impl) parseEdition() (ast.Edition, error) {
	first := p.curr()

	if peek := p.peek(); peek.Kind != token.KindEqual {
		return ast.Edition{}, gotUnexpected(peek, token.KindEqual)
	}
	p.nextToken()

	lit, err := p.parseString()

	if err != nil {
		return ast.Edition{}, err
	}

	if peek := p.peek(); peek.Kind != token.KindSemicolon {
		return ast.Edition{}, gotUnexpected(&p.tokens[p.idx], token.KindSemicolon)
	}

	last := p.nextToken()
	id := p.fm.Merge(token.KindEdition, first.ID, last.ID)
	return ast.Edition{ID: id, Value: lit}, nil
}
