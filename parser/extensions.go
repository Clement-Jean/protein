package parser

import (
	"github.com/Clement-Jean/protein/ast"
	internal_bytes "github.com/Clement-Jean/protein/internal/bytes"
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

		last := p.curr()
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

func (p *impl) parseExtend() (extend ast.Extend, errs []error) {
	first := p.curr()
	id, err := p.parseFullyQualifiedIdentifier()

	if err != nil {
		return ast.Extend{}, []error{err}
	}

	if peek := p.peek(); peek.Kind != token.KindLeftBrace {
		return ast.Extend{}, []error{gotUnexpected(peek, token.KindLeftBrace)}
	}
	p.nextToken()

	peek := p.peek()
	for ; peek.Kind != token.KindRightBrace && peek.Kind != token.KindEOF; peek = p.peek() {
		if peek.Kind == token.KindSemicolon {
			p.nextToken()
			continue
		}

		err = nil
		kind := peek.Kind

		if literal := p.fm.Lookup(peek.ID); literal != nil {
			if k, ok := literalToKind[internal_bytes.ToString(literal)]; ok {
				kind = k
			}
		}

		switch kind {
		case token.KindOption:
			p.nextToken()
			opt, innerErrs := p.parseOption()
			if len(innerErrs) == 0 {
				extend.Options = append(extend.Options, opt)
			}
			errs = append(errs, innerErrs...)
		case token.KindIdentifier:
			var field ast.Field

			if field, err = p.parseField(); err == nil {
				extend.Fields = append(extend.Fields, field)
			}
		default:
			err = gotUnexpected(peek, token.KindOption, token.KindField)
		}

		if err != nil {
			errs = append(errs, err)
			p.advanceTo(exprEnd)

			if p.curr().Kind == token.KindRightBrace {
				extend.Name = id
				extend.ID = p.fm.Merge(token.KindExtend, first.ID, p.curr().ID)
				return extend, errs
			}
		}
	}

	if peek.Kind != token.KindRightBrace {
		errs = append(errs, gotUnexpected(peek, token.KindRightBrace))
		return ast.Extend{}, errs
	}

	last := p.nextToken()
	extend.Name = id
	extend.ID = p.fm.Merge(token.KindExtend, first.ID, last.ID)
	return extend, errs
}
