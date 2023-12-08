package parser

import (
	"bytes"

	"github.com/Clement-Jean/protein/ast"
	internal_bytes "github.com/Clement-Jean/protein/internal/bytes"
	"github.com/Clement-Jean/protein/token"
)

func (p *impl) parseImport() (ast.Import, error) {
	first := p.curr()

	isPublic := false
	isWeak := false
	if peek := p.peek(); peek.Kind == token.KindIdentifier {
		pu := internal_bytes.FromString("public")
		we := internal_bytes.FromString("weak")
		literal := p.fm.Lookup(peek.ID)

		if public := bytes.Compare(literal, pu) == 0; public || bytes.Compare(literal, we) == 0 {
			isPublic = public
			isWeak = !public
		} else {
			return ast.Import{}, gotUnexpected(peek, token.KindPublic, token.KindWeak)
		}
		p.nextToken()
	}

	lit, err := p.parseString()

	if err != nil {
		return ast.Import{}, err
	}

	if peek := p.peek(); peek.Kind != token.KindSemicolon {
		return ast.Import{}, gotUnexpected(peek, token.KindSemicolon)
	}

	last := p.nextToken()
	id := p.fm.Merge(token.KindImport, first.ID, last.ID)
	return ast.Import{ID: id, Value: lit, IsWeak: isWeak, IsPublic: isPublic}, nil
}
