package parser

import (
	"bytes"

	"github.com/Clement-Jean/protein/ast"
	internal_bytes "github.com/Clement-Jean/protein/internal/bytes"
	"github.com/Clement-Jean/protein/token"
)

func (p *impl) parseInt() (ast.Integer, error) {
	if peek := p.peek(); peek.Kind != token.KindInt {
		return ast.Integer{}, gotUnexpected(peek, token.KindInt)
	}

	next := p.nextToken()
	return ast.Integer{ID: next.ID}, nil
}

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

func (p *impl) parseTextFieldName() (ast.Identifier, error) {
	peek := p.peek()

	if peek.Kind == token.KindLeftSquare {
		p.nextToken()
		first, err := p.parseFullyQualifiedIdentifier()

		if err != nil {
			return ast.Identifier{}, err
		}

		parts := []token.UniqueID{first.ID}
		peek := p.peek()
		for ; peek.Kind == token.KindSlash; peek = p.peek() {
			p.nextToken()
			next, err := p.parseFullyQualifiedIdentifier()

			if err != nil {
				return ast.Identifier{}, err
			}

			parts = append(parts, next.ID)
		}

		if peek.Kind != token.KindRightSquare {
			return ast.Identifier{}, gotUnexpected(peek, token.KindRightSquare)
		}

		if len(parts) > 1 {
			id := p.fm.Merge(token.KindIdentifier, parts...)
			return ast.Identifier{ID: id, Parts: parts}, nil
		}

		return ast.Identifier{ID: first.ID}, nil
	} else if name, err := p.parseIdentifier(); err == nil {
		return name, nil
	}

	return ast.Identifier{}, gotUnexpected(peek, token.KindIdentifier, token.KindLeftSquare)
}

func (p *impl) parseTextMessageList(recurseDepth uint8) (list ast.TextMessageList, errs []error) {
	first := p.curr()
	peek := p.peek()
	msgs := []ast.TextMessage{}
	for ; peek.Kind != token.KindRightSquare; peek = p.peek() {
		if peek.Kind == token.KindComma {
			p.nextToken()
		}
		msg, innerErrs := p.parseTextMessage(recurseDepth)

		if len(innerErrs) != 0 {
			errs = append(errs, innerErrs...)

			if p.curr().Kind == token.KindRightSquare {
				id := p.fm.Merge(token.KindTextMessageList, first.ID, p.curr().ID)
				return ast.TextMessageList{ID: id, Values: msgs}, errs
			}

			continue
		}

		msgs = append(msgs, msg)
	}

	if peek.Kind != token.KindRightSquare {
		errs = append(errs, gotUnexpected(peek, token.KindRightSquare))
		return ast.TextMessageList{}, errs
	}

	last := p.nextToken()
	id := p.fm.Merge(token.KindTextMessageList, first.ID, last.ID)
	return ast.TextMessageList{ID: id, Values: msgs}, errs
}

func (p *impl) parseTextScalarList(recurseDepth uint8) (list ast.TextScalarList, errs []error) {
	first := p.curr()
	peek := p.peek()
	values := []ast.Expression{}
	for ; peek.Kind != token.KindRightSquare; peek = p.peek() {
		if peek.Kind == token.KindComma {
			p.nextToken()
		}

		value, innerErrs := p.parseConstant(recurseDepth)

		if len(innerErrs) != 0 {
			errs = append(errs, innerErrs...)
			continue
		}

		values = append(values, value)
	}

	if peek.Kind != token.KindRightSquare {
		errs = append(errs, gotUnexpected(peek, token.KindRightSquare))
		return ast.TextScalarList{}, errs
	}

	last := p.nextToken()
	id := p.fm.Merge(token.KindTextScalarList, first.ID, last.ID)
	return ast.TextScalarList{ID: id, Values: values}, nil
}

func (p *impl) parseTextField(recurseDepth uint8) (field ast.TextField, errs []error) {
	name, err := p.parseTextFieldName()

	if err != nil {
		return ast.TextField{}, []error{err}
	}

	var innerErrs []error
	var value ast.Expression

	peek := p.peek()
	hasColon := peek.Kind == token.KindColon
	if hasColon {
		p.nextToken()
		peek = p.peek()
	}

	if peek.Kind == token.KindLeftSquare {
		curr := p.nextToken()
		peek := p.peek()

		if peek.Kind == token.KindLeftBrace || peek.Kind == token.KindLeftAngle {
			value, innerErrs = p.parseTextMessageList(recurseDepth)
			errs = append(errs, innerErrs...)
		} else if hasColon {
			value, innerErrs = p.parseTextScalarList(recurseDepth)
			errs = append(errs, innerErrs...)
		} else {
			err = gotUnexpected(curr, token.KindColon)
		}

		if err != nil {
			errs = append(errs, err)
			p.advanceTo(exprEnd)
		}
	} else {
		if !hasColon && peek.Kind != token.KindLeftBrace && peek.Kind != token.KindLeftAngle {
			errs = append(errs, gotUnexpected(peek, token.KindColon))
			return ast.TextField{}, errs
		}

		value, innerErrs = p.parseConstant(recurseDepth + 1)
		errs = append(errs, innerErrs...)
	}

	last := p.curr()
	id := p.fm.Merge(token.KindTextField, name.ID, last.ID)
	return ast.TextField{ID: id, Name: name, Value: value}, errs
}

func (p *impl) parseTextMessage(recurseDepth uint8) (msg ast.TextMessage, errs []error) {
	open := p.peek()

	var closeKind token.Kind

	if open.Kind == token.KindLeftBrace {
		closeKind = token.KindRightBrace
	} else if open.Kind == token.KindLeftAngle {
		closeKind = token.KindRightAngle
	} else {
		p.advanceTo(exprEnd)
		return ast.TextMessage{}, []error{gotUnexpected(open, token.KindLeftBrace, token.KindLeftAngle)}
	}

	p.nextToken()
	peek := p.peek()

	for ; peek.Kind != closeKind && peek.Kind != token.KindEOF; peek = p.peek() {
		if peek.Kind == token.KindComma || peek.Kind == token.KindSemicolon {
			p.nextToken()
			continue
		}

		field, innerErrs := p.parseTextField(recurseDepth)

		if len(innerErrs) != 0 {
			errs = append(errs, innerErrs...)
			p.advanceTo(exprEnd)

			if p.curr().Kind == closeKind {
				msg.ID = p.fm.Merge(token.KindTextMessage, open.ID, p.curr().ID)
				return msg, errs
			}

			continue
		}

		msg.Fields = append(msg.Fields, field)
	}

	if peek.Kind != closeKind {
		errs = append(errs, gotUnexpected(peek, closeKind))
		return ast.TextMessage{}, errs
	}
	last := p.nextToken()
	msg.ID = p.fm.Merge(token.KindTextMessage, open.ID, last.ID)
	return msg, errs
}

var expectedConstants = []token.Kind{
	token.KindInt,
	token.KindFloat,
	token.KindIdentifier,
	token.KindStr,
	token.KindLeftBrace,
	token.KindLeftAngle,
}

func (p *impl) parseConstant(recurseDepth uint8) (value ast.Expression, errs []error) {
	curr := p.curr()

	if recurseDepth > 30 { // TODO make it configurable
		return nil, []error{&Error{
			ID:  curr.ID,
			Msg: "Too many nested constants",
		}}
	}

	var err error

	peek := p.peek()
	switch peek.Kind {
	case token.KindInt:
		value = ast.Integer{ID: peek.ID}
	case token.KindFloat:
		value = ast.Float{ID: peek.ID}
	case token.KindIdentifier:
		literal := p.fm.Lookup(peek.ID)
		tr := internal_bytes.FromString("true")
		fa := internal_bytes.FromString("false")
		if t := bytes.Compare(literal, tr) == 0; t || bytes.Compare(literal, fa) == 0 {
			value = ast.Boolean{ID: peek.ID}
		} else {
			value = ast.Identifier{ID: peek.ID}
		}
	case token.KindStr:
		value = ast.String{ID: peek.ID}
	case token.KindLeftBrace:
		var innerErrs []error
		value, innerErrs = p.parseTextMessage(recurseDepth)
		errs = append(errs, innerErrs...)

		if p.curr().Kind == token.KindRightBrace {
			return value, errs
		}
	case token.KindLeftAngle:
		var innerErrs []error
		value, innerErrs = p.parseTextMessage(recurseDepth)
		errs = append(errs, innerErrs...)

		if p.curr().Kind == token.KindRightBrace {
			return value, errs
		}
	default:
		err = gotUnexpected(peek, expectedConstants...)
	}

	if err != nil {
		errs = append(errs, err)
		p.advanceTo(exprEnd)
		return value, errs
	}

	p.nextToken()
	return value, errs
}
