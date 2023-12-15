package parser

import (
	"github.com/Clement-Jean/protein/ast"
	internal_bytes "github.com/Clement-Jean/protein/internal/bytes"
	"github.com/Clement-Jean/protein/token"
)

func (p *impl) parseOneof() (oneof ast.Oneof, err error) {
	first := p.curr()

	if peek := p.peek(); peek.Kind != token.KindIdentifier {
		return ast.Oneof{}, gotUnexpected(peek, token.KindIdentifier)
	}
	p.nextToken()

	name, err := p.parseIdentifier()

	if err != nil {
		return ast.Oneof{}, err
	}

	if peek := p.peek(); peek.Kind != token.KindLeftBrace {
		return ast.Oneof{}, gotUnexpected(peek, token.KindLeftBrace)
	}
	p.nextToken()

	peek := p.peek()
	for ; peek.Kind != token.KindRightBrace && peek.Kind != token.KindEOF; peek = p.peek() {
		if peek.Kind == token.KindSemicolon {
			p.nextToken()
			continue
		}

		kind := peek.Kind

		if literal := p.fm.Lookup(peek.ID); literal != nil {
			if k, ok := literalToKind[internal_bytes.ToString(literal)]; ok {
				kind = k
			}
		}

		switch kind {
		case token.KindOption:
			var option ast.Option

			p.nextToken() // point to option keyword
			if option, err = p.parseOption(); err == nil {
				oneof.Options = append(oneof.Options, option)
			}
		case token.KindIdentifier:
			var field ast.Field

			if field, err = p.parseField(); err == nil {
				oneof.Fields = append(oneof.Fields, field)
			}
		default:
			err = gotUnexpected(peek, token.KindOption, token.KindIdentifier, token.KindRightBrace)
		}

		if err != nil {
			// TODO report error
			// TODO p.advanceTo(exprEnd)
			return ast.Oneof{}, err
		}
	}

	if peek.Kind != token.KindRightBrace {
		return ast.Oneof{}, gotUnexpected(peek, token.KindRightBrace)
	}

	last := p.nextToken()

	oneof.ID = p.fm.Merge(token.KindOneOf, first.ID, last.ID)
	oneof.Name = name
	return oneof, nil
}
