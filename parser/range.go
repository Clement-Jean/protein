package parser

import (
	"bytes"

	"github.com/Clement-Jean/protein/ast"
	"github.com/Clement-Jean/protein/token"
)

func (p *impl) parseRange() (ast.Range, error) {
	first, err := p.parseInt()

	if err != nil {
		return ast.Range{}, err
	}

	last := first
	if peek := p.peek(); peek.Kind == token.KindIdentifier && bytes.Compare(p.fm.Lookup(peek.ID), []byte("to")) == 0 {
		p.nextToken()

		switch peek = p.peek(); peek.Kind {
		case token.KindInt:
			last, _ = p.parseInt()
		case token.KindIdentifier:
			if bytes.Compare(p.fm.Lookup(peek.ID), []byte("max")) == 0 {
				id, _ := p.parseIdentifier()
				last = ast.Integer{ID: id.ID}
				break
			}
			fallthrough
		default:
			return ast.Range{}, gotUnexpected(peek, token.KindMax, token.KindInt)
		}
	}

	var id token.UniqueID

	if last != first {
		id = p.fm.Merge(token.KindRange, first.ID, last.ID)
	} else {
		id = first.ID
	}
	return ast.Range{ID: id, Start: first, End: last}, err
}
