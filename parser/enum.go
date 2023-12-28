package parser

import (
	"github.com/Clement-Jean/protein/ast"
	"github.com/Clement-Jean/protein/internal/bytes"
	"github.com/Clement-Jean/protein/token"
)

func (p *impl) parseEnumValue() (value ast.EnumValue, errs []error) {
	name, _ := p.parseIdentifier()

	if peek := p.peek(); peek.Kind != token.KindEqual {
		return ast.EnumValue{}, []error{gotUnexpected(peek, token.KindEqual)}
	}
	p.nextToken()

	tag, err := p.parseInt()

	if err != nil {
		return ast.EnumValue{}, []error{err}
	}

	var options []ast.Option
	var innerErrs []error
	var optionsID token.UniqueID

	peek := p.peek()
	if peek.Kind == token.KindLeftSquare {
		firstOption := p.nextToken()
		options, innerErrs = p.parseInlineOptions()
		errs = append(errs, innerErrs...)

		var lastOptionID token.UniqueID
		if len(options) != 0 {
			lastOptionID = options[len(options)-1].ID
		} else {
			lastOptionID = firstOption.ID
		}
		optionsID = p.fm.Merge(token.KindOption, firstOption.ID, lastOptionID)
	}

	if peek = p.peek(); peek.Kind != token.KindSemicolon {
		errs = append(errs, gotUnexpected(peek, token.KindSemicolon))
		return ast.EnumValue{}, errs
	}

	last := p.nextToken()
	id := p.fm.Merge(token.KindEnumValue, name.ID, last.ID)
	value.ID = id
	value.Name = name
	value.Tag = tag
	value.Options = options
	value.OptionsID = optionsID
	return value, errs
}

func (p *impl) parseEnum() (enum ast.Enum, errs []error) {
	first := p.curr()
	id, err := p.parseIdentifier()

	if err != nil {
		return ast.Enum{}, []error{err}
	}

	if peek := p.peek(); peek.Kind != token.KindLeftBrace {
		return ast.Enum{}, []error{gotUnexpected(peek, token.KindLeftBrace)}
	}
	p.nextToken()

	peek := p.peek()
	for ; peek.Kind != token.KindRightBrace && peek.Kind != token.KindEOF; peek = p.peek() {
		if peek.Kind == token.KindSemicolon {
			p.nextToken()
			continue
		}

		err = nil
		prevErrsLen := len(errs)
		kind := peek.Kind

		if literal := p.fm.Lookup(peek.ID); literal != nil {
			if k, ok := literalToKind[bytes.ToString(literal)]; ok {
				kind = k
			}
		}

		switch kind {
		case token.KindOption:
			p.nextToken() // point to option keyword
			opt, innerErrs := p.parseOption()
			if len(innerErrs) == 0 {
				enum.Options = append(enum.Options, opt)
			}
			errs = append(errs, innerErrs...)
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
			value, innerErrs := p.parseEnumValue()
			if len(innerErrs) == 0 {
				enum.Values = append(enum.Values, value)
			}
			errs = append(errs, innerErrs...)
		default:
			err = gotUnexpected(peek, token.KindOption, token.KindReserved, token.KindIdentifier)
		}

		if err != nil || prevErrsLen != len(errs) {
			if err != nil {
				errs = append(errs, err)
			}
			p.advanceTo(exprEnd)

			if p.curr().Kind == token.KindRightBrace {
				enum.Name = id
				enum.ID = p.fm.Merge(token.KindEnum, first.ID, p.curr().ID)
				return enum, errs
			}
		}
	}

	if peek.Kind != token.KindRightBrace {
		errs = append(errs, gotUnexpected(peek, token.KindRightBrace))
		return ast.Enum{}, errs
	}

	last := p.nextToken()
	enum.Name = id
	enum.ID = p.fm.Merge(token.KindEnum, first.ID, last.ID)
	return enum, errs
}
