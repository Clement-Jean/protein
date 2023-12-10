package parser

import (
	"bytes"

	"github.com/Clement-Jean/protein/ast"
	internal_bytes "github.com/Clement-Jean/protein/internal/bytes"
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

var expectedConstants = []token.Kind{
	token.KindInt,
	token.KindFloat,
	token.KindIdentifier,
	token.KindStr,
	token.KindLeftBrace,
	token.KindLeftAngle,
}

func (p *impl) parseConstant(recurseDepth uint8) (ast.Expression, error) {
	curr := p.curr()

	if recurseDepth > 30 { // TODO make it configurable
		return nil, &Error{
			ID:  curr.ID,
			Msg: "Too many nested constants",
		}
	}

	switch peek := p.peek(); peek.Kind {
	case token.KindInt:
	case token.KindFloat:
	case token.KindIdentifier:
	case token.KindStr:
	case token.KindLeftBrace:
	case token.KindLeftAngle:
	default:
		return nil, gotUnexpected(peek, expectedConstants...)
	}

	next := p.nextToken()

	switch next.Kind {
	case token.KindInt:
		return &ast.Integer{ID: next.ID}, nil
	case token.KindFloat:
		return &ast.Float{ID: next.ID}, nil
	case token.KindIdentifier:
		literal := p.fm.Lookup(next.ID)
		tr := internal_bytes.FromString("true")
		fa := internal_bytes.FromString("false")
		if t := bytes.Compare(literal, tr) == 0; t || bytes.Compare(literal, fa) == 0 {
			return &ast.Boolean{ID: next.ID}, nil
		}
		return &ast.Identifier{ID: next.ID}, nil
	case token.KindStr:
		return &ast.String{ID: next.ID}, nil
	case token.KindLeftBrace:
		panic("to implement")
	case token.KindLeftAngle:
		panic("to implement")
	default:
		panic("unreachable")
	}
}
