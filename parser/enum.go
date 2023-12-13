package parser

import (
	"github.com/Clement-Jean/protein/ast"
	"github.com/Clement-Jean/protein/internal/bytes"
	"github.com/Clement-Jean/protein/token"
)

func (p *impl) parseEnumValue() (ast.EnumValue, error) {
	name, _ := p.parseIdentifier()

	if peek := p.peek(); peek.Kind != token.KindEqual {
		return ast.EnumValue{}, gotUnexpected(peek, token.KindEqual)
	}
	p.nextToken()

	tag, err := p.parseInt()

	if err != nil {
		return ast.EnumValue{}, err
	}

	var options []ast.Option
	var optionsID token.UniqueID
	var firstOption *token.Token
	var lastOption *token.Token

	peek := p.peek()
	if peek.Kind == token.KindLeftSquare {
		firstOption = p.nextToken()
		options, err = p.parseInlineOptions()

		if err != nil {
			return ast.EnumValue{}, err
		}

		lastOption = p.curr()
		p.nextToken()
		optionsID = p.fm.Merge(token.KindOption, firstOption.ID, lastOption.ID)
	}

	if peek = p.peek(); peek.Kind != token.KindSemicolon {
		return ast.EnumValue{}, gotUnexpected(peek, token.KindSemicolon)
	}

	last := p.nextToken()
	id := p.fm.Merge(token.KindEnumValue, name.ID, last.ID)
	return ast.EnumValue{ID: id, Name: name, Tag: tag, Options: options, OptionsID: optionsID}, nil
}

func (p *impl) parseEnum() (ast.Enum, error) {
	first := p.curr()
	id, err := p.parseIdentifier()

	if err != nil {
		return ast.Enum{}, err
	}

	if peek := p.peek(); peek.Kind != token.KindLeftBrace {
		return ast.Enum{}, gotUnexpected(peek, token.KindLeftBrace)
	}
	p.nextToken()

	var enum ast.Enum

	peek := p.peek()
	for ; peek.Kind != token.KindRightBrace && peek.Kind != token.KindEOF; peek = p.peek() {
		if peek.Kind == token.KindSemicolon {
			p.nextToken()
			continue
		}

		kind := peek.Kind

		if literal := p.fm.Lookup(peek.ID); literal != nil {
			if k, ok := literalToKind[bytes.ToString(literal)]; ok {
				kind = k
			}
		}

		switch kind {
		case token.KindOption:
			var option ast.Option

			p.nextToken() // point to option keyword
			if option, err = p.parseOption(); err == nil {
				enum.Options = append(enum.Options, option)
			}
		case token.KindReserved:
			p.nextToken() // point to reserved keyword

			peek := p.peek()
			if peek.Kind == token.KindInt {
				var reserved ast.ReservedTags

				if reserved, err = p.parseReservedTags(); err == nil {
					enum.ReservedTags = append(enum.ReservedTags, reserved)
				}
			} else if peek.Kind == token.KindStr {
				var reserved ast.ReservedNames

				if reserved, err = p.parseReservedNames(); err == nil {
					enum.ReservedNames = append(enum.ReservedNames, reserved)
				}
			}
		case token.KindIdentifier:
			var value ast.EnumValue

			if value, err = p.parseEnumValue(); err == nil {
				enum.Values = append(enum.Values, value)
			}
		default:
			err = gotUnexpected(peek, token.KindOption, token.KindReserved, token.KindIdentifier, token.KindRightBrace)
		}

		if err != nil {
			// TODO report error
			// TODO p.advanceTo(exprEnd)
			return ast.Enum{}, err
		}
	}

	if peek.Kind != token.KindRightBrace {
		return ast.Enum{}, gotUnexpected(peek, token.KindRightBrace)
	}

	last := p.nextToken()
	enum.Name = id
	enum.ID = p.fm.Merge(token.KindEnum, first.ID, last.ID)
	return enum, nil
}
