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

func (p *impl) parseOptionIdentifierEqualValue() (name ast.Identifier, value ast.Expression, errs []error) {
	var err error

	name, err = p.parseOptionName()

	if err != nil {
		p.advanceTo(exprEnd)
		return name, nil, []error{err}
	}

	if peek := p.peek(); peek.Kind != token.KindEqual {
		p.advanceTo(exprEnd)
		errs = append(errs, gotUnexpected(peek, token.KindEqual))
		return name, value, errs
	}
	p.nextToken()

	value, errs = p.parseConstant(1)

	if len(errs) != 0 {
		return name, value, errs
	}

	return name, value, errs
}

func (p *impl) parseInlineOption() (opt ast.Option, errs []error) {
	opt.Name, opt.Value, errs = p.parseOptionIdentifierEqualValue()

	if errs != nil {
		return ast.Option{}, errs
	}

	id := p.fm.Merge(token.KindOption, opt.Name.ID, opt.Value.GetID())
	opt.ID = id
	return opt, errs
}

func (p *impl) parseInlineOptions() (opts []ast.Option, errs []error) {
	peek := p.peek()
	curr := p.curr()

	for ; curr.Kind != token.KindRightSquare && curr.Kind != token.KindEOF && peek.Kind != token.KindRightSquare && peek.Kind != token.KindEOF; peek = p.peek() {
		if peek.Kind == token.KindComma {
			p.nextToken()
		}

		option, innerErrs := p.parseInlineOption()

		if len(innerErrs) != 0 {
			errs = append(errs, innerErrs...)
			continue
		}

		opts = append(opts, option)
		curr = p.curr()
	}

	if curr.Kind == token.KindRightSquare {
		return opts, errs
	}

	if peek.Kind != token.KindRightSquare {
		errs = append(errs, gotUnexpected(peek, token.KindRightSquare))
		return nil, errs
	}
	p.nextToken()
	return opts, errs
}

func (p *impl) parseOption() (opt ast.Option, errs []error) {
	first := p.curr()
	name, value, innerErrs := p.parseOptionIdentifierEqualValue()

	if len(innerErrs) != 0 {
		errs = append(errs, innerErrs...)

		if p.curr().Kind == token.KindSemicolon {
			id := p.fm.Merge(token.KindOption, first.ID, p.curr().ID)
			return ast.Option{ID: id, Name: name, Value: value}, errs
		}
	}

	if peek := p.peek(); peek.Kind != token.KindSemicolon {
		errs = append(errs, gotUnexpected(peek, token.KindSemicolon))
		return ast.Option{}, errs
	}

	last := p.nextToken()
	id := p.fm.Merge(token.KindOption, first.ID, last.ID)
	return ast.Option{ID: id, Name: name, Value: value}, errs
}
