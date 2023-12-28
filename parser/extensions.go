package parser

import (
	"github.com/Clement-Jean/protein/ast"
	internal_bytes "github.com/Clement-Jean/protein/internal/bytes"
	"github.com/Clement-Jean/protein/token"
)

func (p *impl) parseExtensionRange() (er ast.ExtensionRange, errs []error) {
	first := p.curr()
	r, err := p.parseRange()

	if err != nil {
		return ast.ExtensionRange{}, []error{err}
	}

	ranges := []ast.Range{r}
	peek := p.peek()
	for ; peek.Kind == token.KindComma; peek = p.peek() {
		p.nextToken()

		if r, err = p.parseRange(); err != nil {
			return er, []error{err}
		}

		ranges = append(ranges, r)
	}

	var opts []ast.Option
	var innerErrs []error
	var optsID token.UniqueID

	if peek := p.peek(); peek.Kind == token.KindLeftSquare {
		first := p.nextToken()
		opts, innerErrs = p.parseInlineOptions()
		errs = append(errs, innerErrs...)
		last := p.curr()
		optsID = p.fm.Merge(token.KindOption, first.ID, last.ID)
	}

	if peek := p.peek(); peek.Kind != token.KindSemicolon {
		errs = append(errs, gotUnexpected(peek, token.KindSemicolon))
		return ast.ExtensionRange{}, errs
	}
	last := p.nextToken()

	er.Options = opts
	er.OptionsID = optsID
	er.ID = p.fm.Merge(token.KindExtensions, first.ID, last.ID)
	er.Ranges = ranges
	return er, errs
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
			field, innerErrs := p.parseField()
			if len(innerErrs) == 0 {
				extend.Fields = append(extend.Fields, field)
			}
			errs = append(errs, innerErrs...)
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
