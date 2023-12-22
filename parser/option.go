package parser

import (
	"github.com/Clement-Jean/protein/ast"
	"github.com/Clement-Jean/protein/token"
)

func (p *impl) parseOptionName() (ast.Identifier, error) {
	var ids []token.UniqueID

	if peek := p.peek(); peek.Kind == token.KindLeftParen {
		paren := p.nextToken()
		ids = append(ids, paren.ID)

		if peek := p.peek(); peek.Kind == token.KindDot {
			p.nextToken()
		}

		name, err := p.parseFullyQualifiedIdentifier()

		if err != nil {
			return name, err
		}

		ids = append(ids, name.ID)

		if peek := p.peek(); peek.Kind != token.KindRightParen {
			return ast.Identifier{}, gotUnexpected(peek, token.KindLeftParen)
		}
		paren = p.nextToken()
		ids = append(ids, paren.ID)

		if peek := p.peek(); peek.Kind == token.KindDot {
			p.nextToken()
		}

		name, err = p.parseFullyQualifiedIdentifier()

		if err == nil {
			ids = append(ids, name.ID)
		}
	} else {
		name, err := p.parseFullyQualifiedIdentifier()

		if err != nil {
			return name, err
		}

		ids = append(ids, name.ID)
	}

	if len(ids) > 1 {
		id := p.fm.Merge(token.KindIdentifier, ids...)
		return ast.Identifier{ID: id, Parts: ids}, nil
	}

	//TODO check that len(ids) == 1 otherwise error

	return ast.Identifier{ID: ids[0]}, nil
}

func (p *impl) parseOptionValue() (ast.Expression, error) {
	lit, err := p.parseConstant(1)

	if err != nil {
		return nil, err
	}

	return lit, err
}

func (p *impl) parseOptionIdentifierEqualValue() (ast.Identifier, ast.Expression, error) {
	name, err := p.parseOptionName()

	if err != nil {
		return name, nil, err
	}

	if peek := p.peek(); peek.Kind != token.KindEqual {
		return name, nil, gotUnexpected(peek, token.KindEqual)
	}
	p.nextToken()

	value, err := p.parseOptionValue()

	if err != nil {
		return name, nil, err
	}

	return name, value, nil
}

func (p *impl) parseInlineOption() (ast.Option, error) {
	name, value, err := p.parseOptionIdentifierEqualValue()

	if err != nil {
		return ast.Option{}, err
	}

	id := p.fm.Merge(token.KindOption, name.ID, value.GetID())
	return ast.Option{ID: id, Name: name, Value: value}, nil
}

func (p *impl) parseInlineOptions() ([]ast.Option, error) {
	var options []ast.Option
	peek := p.peek()
	curr := p.curr()

	for ; curr.Kind != token.KindRightSquare && curr.Kind != token.KindEOF && peek.Kind != token.KindRightSquare && peek.Kind != token.KindEOF; peek = p.peek() {
		if peek.Kind == token.KindComma {
			p.nextToken()
		}

		option, err := p.parseInlineOption()

		if err != nil {
			return nil, err
		}

		options = append(options, option)
		curr = p.curr()
	}

	if curr.Kind == token.KindRightSquare {
		return options, nil
	}

	if peek.Kind != token.KindRightSquare {
		return nil, gotUnexpected(peek, token.KindRightSquare)
	}
	p.nextToken()
	return options, nil
}

func (p *impl) parseOption() (ast.Option, error) {
	first := p.curr()
	name, value, err := p.parseOptionIdentifierEqualValue()

	if err != nil {
		return ast.Option{}, err
	}

	if peek := p.peek(); peek.Kind != token.KindSemicolon {
		return ast.Option{}, gotUnexpected(peek, token.KindSemicolon)
	}

	last := p.nextToken()
	id := p.fm.Merge(token.KindOption, first.ID, last.ID)
	return ast.Option{ID: id, Name: name, Value: value}, nil
}
