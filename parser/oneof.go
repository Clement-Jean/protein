package parser

import (
	"github.com/Clement-Jean/protein/ast"
	internal_bytes "github.com/Clement-Jean/protein/internal/bytes"
	"github.com/Clement-Jean/protein/token"
)

func (p *impl) parseOneof() (oneof ast.Oneof, errs []error) {
	first := p.curr()

	if peek := p.peek(); peek.Kind != token.KindIdentifier {
		return ast.Oneof{}, []error{gotUnexpected(peek, token.KindIdentifier)}
	}
	p.nextToken()

	name, err := p.parseIdentifier()

	if err != nil {
		return ast.Oneof{}, []error{err}
	}

	if peek := p.peek(); peek.Kind != token.KindLeftBrace {
		return ast.Oneof{}, []error{gotUnexpected(peek, token.KindLeftBrace)}
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
			p.nextToken() // point to option keyword
			opt, innerErrs := p.parseOption()
			if len(innerErrs) == 0 {
				oneof.Options = append(oneof.Options, opt)
			}
			errs = append(errs, innerErrs...)
		case token.KindIdentifier:
			field, innerErrs := p.parseField()
			if len(innerErrs) == 0 {
				oneof.Fields = append(oneof.Fields, field)
			}
			errs = append(errs, innerErrs...)
		default:
			err = gotUnexpected(peek, token.KindOption, token.KindIdentifier)
		}

		if err != nil {
			errs = append(errs, err)
			p.advanceTo(exprEnd)

			if p.curr().Kind == token.KindRightBrace {
				oneof.Name = name
				oneof.ID = p.fm.Merge(token.KindOneOf, first.ID, p.curr().ID)
				return oneof, errs
			}
		}
	}

	if peek.Kind != token.KindRightBrace {
		errs = append(errs, gotUnexpected(peek, token.KindRightBrace))
		return ast.Oneof{}, errs
	}

	last := p.nextToken()

	oneof.ID = p.fm.Merge(token.KindOneOf, first.ID, last.ID)
	oneof.Name = name
	return oneof, errs
}
