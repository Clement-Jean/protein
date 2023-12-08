package parser

import (
	"github.com/Clement-Jean/protein/ast"
	"github.com/Clement-Jean/protein/token"
)

func (p *impl) parseSyntax() (ast.Syntax, error) {
	first := p.curr()

	if peek := p.peek(); peek.Kind != token.KindEqual {
		return ast.Syntax{}, gotUnexpected(peek, token.KindEqual)
	}
	p.nextToken()

	lit, err := p.parseString()

	if err != nil {
		return ast.Syntax{}, err
	}

	if peek := p.peek(); peek.Kind != token.KindSemicolon {
		return ast.Syntax{}, gotUnexpected(&p.tokens[p.idx], token.KindSemicolon)
	}

	last := p.nextToken()
	id := p.fm.Merge(token.KindSyntax, first.ID, last.ID)
	return ast.Syntax{ID: id, Value: lit}, nil
}
