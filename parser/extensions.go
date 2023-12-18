package parser

import (
	"github.com/Clement-Jean/protein/ast"
	"github.com/Clement-Jean/protein/token"
)

func (p *impl) parseExtensionRange() (er ast.ExtensionRange, err error) {
	first := p.curr()
	r, err := p.parseRange()

	if err != nil {
		return ast.ExtensionRange{}, err
	}

	ranges := []ast.Range{r}
	peek := p.peek()
	for ; peek.Kind == token.KindComma; peek = p.peek() {
		p.nextToken()

		if r, err = p.parseRange(); err != nil {
			return er, err
		}

		ranges = append(ranges, r)
	}

	var opts []ast.Option
	var optsID token.UniqueID
	if peek := p.peek(); peek.Kind == token.KindLeftSquare {
		first := p.nextToken()
		opts, err = p.parseInlineOptions()

		if err != nil {
			return ast.ExtensionRange{}, err
		}

		last := p.nextToken()
		optsID = p.fm.Merge(token.KindOption, first.ID, last.ID)
	}

	if peek := p.peek(); peek.Kind != token.KindSemicolon {
		return ast.ExtensionRange{}, gotUnexpected(peek, token.KindSemicolon)
	}
	last := p.nextToken()

	er.Options = opts
	er.OptionsID = optsID
	er.ID = p.fm.Merge(token.KindExtensions, first.ID, last.ID)
	er.Ranges = ranges
	return er, err
}
