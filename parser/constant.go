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

func (p *impl) parseTextMessageList(recurseDepth uint8) (ast.TextMessageList, error) {
	first := p.curr()
	peek := p.peek()
	msgs := []ast.TextMessage{}
	for ; peek.Kind != token.KindRightSquare; peek = p.peek() {
		if peek.Kind == token.KindComma {
			p.nextToken()
		}
		msg, err := p.parseTextMessage(recurseDepth)

		if err != nil {
			return ast.TextMessageList{}, err
		}

		msgs = append(msgs, msg)
	}

	if peek.Kind != token.KindRightSquare {
		return ast.TextMessageList{}, gotUnexpected(peek, token.KindRightSquare)
	}

	last := p.nextToken()
	id := p.fm.Merge(token.KindTextMessageList, first.ID, last.ID)
	return ast.TextMessageList{ID: id, Values: msgs}, nil
}

func (p *impl) parseTextScalarList(recurseDepth uint8) (ast.TextScalarList, error) {
	first := p.curr()
	peek := p.peek()
	values := []ast.Expression{}
	for ; peek.Kind != token.KindRightSquare; peek = p.peek() {
		if peek.Kind == token.KindComma {
			p.nextToken()
		}

		value, err := p.parseConstant(recurseDepth)

		if err != nil {
			return ast.TextScalarList{}, err
		}

		values = append(values, value)
	}

	if peek.Kind != token.KindRightSquare {
		return ast.TextScalarList{}, gotUnexpected(peek, token.KindRightSquare)
	}

	last := p.nextToken()
	id := p.fm.Merge(token.KindTextScalarList, first.ID, last.ID)
	return ast.TextScalarList{ID: id, Values: values}, nil
}

func (p *impl) parseTextField(recurseDepth uint8) (ast.TextField, error) {
	name, err := p.parseTextFieldName()

	if err != nil {
		return ast.TextField{}, err
	}

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
			value, err = p.parseTextMessageList(recurseDepth)
		} else if hasColon {
			value, err = p.parseTextScalarList(recurseDepth)
		} else {
			return ast.TextField{}, gotUnexpected(curr, token.KindColon)
		}

		if err != nil {
			return ast.TextField{}, err
		}
	} else {
		if !hasColon && peek.Kind != token.KindLeftBrace && peek.Kind != token.KindLeftAngle {
			return ast.TextField{}, gotUnexpected(peek, token.KindColon)
		}

		value, err = p.parseConstant(recurseDepth)

		if err != nil {
			return ast.TextField{}, err
		}
	}

	last := p.curr()
	id := p.fm.Merge(token.KindTextField, name.ID, last.ID)
	return ast.TextField{ID: id, Name: name, Value: value}, err
}

func (p *impl) parseTextMessage(recurseDepth uint8) (msg ast.TextMessage, err error) {
	open := p.peek()

	var closeKind token.Kind

	if open.Kind == token.KindLeftBrace {
		closeKind = token.KindRightBrace
	} else if open.Kind == token.KindLeftAngle {
		closeKind = token.KindRightAngle
	} else {
		return ast.TextMessage{}, gotUnexpected(open, token.KindLeftBrace, token.KindLeftAngle)
	}

	p.nextToken()
	peek := p.peek()
	msg = ast.TextMessage{}

	for ; peek.Kind != closeKind && peek.Kind != token.KindEOF; peek = p.peek() {
		if peek.Kind == token.KindComma || peek.Kind == token.KindSemicolon {
			p.nextToken()
			continue
		}

		field, err := p.parseTextField(recurseDepth)

		if err != nil {
			return ast.TextMessage{}, err
		}

		msg.Fields = append(msg.Fields, field)
	}

	if peek.Kind != closeKind {
		return ast.TextMessage{}, gotUnexpected(peek, closeKind)
	}
	last := p.nextToken()

	msg.ID = p.fm.Merge(token.KindTextMessage, open.ID, last.ID)
	return msg, nil
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

	var value ast.Expression
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
		value, err = p.parseTextMessage(recurseDepth)
	case token.KindLeftAngle:
		value, err = p.parseTextMessage(recurseDepth)
	default:
		return nil, gotUnexpected(peek, expectedConstants...)
	}

	if err != nil {
		return nil, err
	}

	p.nextToken()
	return value, err
}
