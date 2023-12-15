package parser

import (
	"github.com/Clement-Jean/protein/ast"
	internal_bytes "github.com/Clement-Jean/protein/internal/bytes"
	"github.com/Clement-Jean/protein/token"
)

var literalToType = map[string]ast.FieldType{
	"double":   ast.FieldTypeDouble,
	"float":    ast.FieldTypeFloat,
	"int64":    ast.FieldTypeInt64,
	"uint64":   ast.FieldTypeUint64,
	"int32":    ast.FieldTypeInt32,
	"fixed64":  ast.FieldTypeFixed64,
	"fixed32":  ast.FieldTypeFixed32,
	"bool":     ast.FieldTypeBool,
	"string":   ast.FieldTypeString,
	"group":    ast.FieldTypeGroup,
	"bytes":    ast.FieldTypeBytes,
	"uint32":   ast.FieldTypeUint32,
	"sfixed32": ast.FieldTypeSfixed32,
	"sfixed64": ast.FieldTypeSfixed64,
	"sint32":   ast.FieldTypeSint32,
	"sint64":   ast.FieldTypeSint64,
}

var literalToLabel = map[string]ast.FieldLabel{
	"optional": ast.FieldLabelOptional,
	"required": ast.FieldLabelRequired,
	"repeated": ast.FieldLabelRepeated,
}

func (p *impl) parseFieldIdentifierTagOption() (ast.Field, error) {
	name, err := p.parseIdentifier()

	if err != nil {
		return ast.Field{}, err
	}

	if peek := p.peek(); peek.Kind != token.KindEqual {
		return ast.Field{}, gotUnexpected(peek, token.KindEqual)
	}
	p.nextToken()

	tag, err := p.parseInt()

	if err != nil {
		return ast.Field{}, err
	}

	var opts []ast.Option
	var optsID token.UniqueID
	if peek := p.peek(); peek.Kind == token.KindLeftSquare {
		first := p.nextToken()
		opts, err = p.parseInlineOptions()

		if err != nil {
			return ast.Field{}, err
		}

		last := p.nextToken()
		optsID = p.fm.Merge(token.KindOption, first.ID, last.ID)
	}

	if peek := p.peek(); peek.Kind != token.KindSemicolon {
		return ast.Field{}, gotUnexpected(peek, token.KindSemicolon)
	}

	return ast.Field{Name: name, Tag: tag, Options: opts, OptionsID: optsID}, nil
}

func (p *impl) parseField() (field ast.Field, err error) {
	id, _ := p.parseFullyQualifiedIdentifier()
	literal := internal_bytes.ToString(p.fm.Lookup(id.ID))

	switch label, ok := literalToLabel[literal]; ok {
	case true:
		field.LabelID = id.ID
		field.Label = label
		id, err = p.parseIdentifier()

		if err != nil {
			return ast.Field{}, err
		}

		literal = internal_bytes.ToString(p.fm.Lookup(id.ID))
		fallthrough
	case false:
		if t, ok := literalToType[literal]; ok {
			field.Type = t
			field.TypeID = id.ID
			break
		}
		field.Type = ast.FieldTypeUnknown // could be an error or message/enum
		field.TypeID = id.ID
	}

	fieldInfo, err := p.parseFieldIdentifierTagOption()

	if err != nil {
		return ast.Field{}, err
	}

	last := p.nextToken()
	field.ID = p.fm.Merge(token.KindField, id.ID, last.ID)
	field.Name = fieldInfo.Name
	field.Tag = fieldInfo.Tag
	field.Options = fieldInfo.Options
	field.OptionsID = fieldInfo.OptionsID
	return field, nil
}

func (p *impl) parseMapField() (field ast.Field, err error) {
	if peek := p.peek(); peek.Kind != token.KindIdentifier {
		return ast.Field{}, gotUnexpected(peek, token.KindIdentifier)
	}
	first := p.nextToken()

	if peek := p.peek(); peek.Kind != token.KindLeftAngle {
		return ast.Field{}, gotUnexpected(peek, token.KindLeftAngle)
	}
	p.nextToken()

	_, err = p.parseIdentifier()

	if err != nil {
		return ast.Field{}, err
	}

	if peek := p.peek(); peek.Kind != token.KindComma {
		return ast.Field{}, gotUnexpected(peek, token.KindComma)
	}
	p.nextToken()

	_, err = p.parseIdentifier()

	if err != nil {
		return ast.Field{}, err
	}

	if peek := p.peek(); peek.Kind != token.KindRightAngle {
		return ast.Field{}, gotUnexpected(peek, token.KindRightAngle)
	}
	endType := p.tokens[p.idx]
	p.nextToken()

	fieldInfo, err := p.parseFieldIdentifierTagOption()

	if err != nil {
		return ast.Field{}, err
	}

	last := p.nextToken()
	field.Name = fieldInfo.Name
	field.Tag = fieldInfo.Tag
	field.Options = fieldInfo.Options
	field.OptionsID = fieldInfo.OptionsID
	field.Type = ast.FieldTypeMessage
	field.TypeID = p.fm.Merge(token.KindMap, first.ID, endType.ID)
	field.ID = p.fm.Merge(token.KindField, first.ID, last.ID)
	return field, nil
}

func (p *impl) parseMessage(recurseDepth uint8) (ast.Message, error) {
	first := p.curr()

	if recurseDepth > 30 { // TODO make it configurable
		return ast.Message{}, &Error{
			ID:  first.ID,
			Msg: "Too many nested messages",
		}
	}

	id, err := p.parseIdentifier()

	if err != nil {
		return ast.Message{}, err
	}

	if peek := p.peek(); peek.Kind != token.KindLeftBrace {
		return ast.Message{}, gotUnexpected(peek, token.KindLeftBrace)
	}
	p.nextToken()

	var msg ast.Message

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
				msg.Options = append(msg.Options, option)
			}
		case token.KindReserved:
			p.nextToken() // point to reserved keyword

			peek := p.peek()
			if peek.Kind == token.KindInt {
				var reserved ast.ReservedTags

				if reserved, err = p.parseReservedTags(); err == nil {
					msg.ReservedTags = append(msg.ReservedTags, reserved)
				}
			} else if peek.Kind == token.KindStr {
				var reserved ast.ReservedNames

				if reserved, err = p.parseReservedNames(); err == nil {
					msg.ReservedNames = append(msg.ReservedNames, reserved)
				}
			}
		case token.KindMap:
			var field ast.Field

			if field, err = p.parseMapField(); err == nil {
				msg.Fields = append(msg.Fields, field)
			}
		case token.KindOneOf:
			var oneof ast.Oneof

			if oneof, err = p.parseOneof(); err == nil {
				msg.Oneofs = append(msg.Oneofs, oneof)
			}
		case token.KindEnum:
			var enum ast.Enum

			p.nextToken() // point to enum keyword
			if enum, err = p.parseEnum(); err == nil {
				msg.Enums = append(msg.Enums, enum)
			}
		case token.KindMessage:
			var innerMsg ast.Message

			p.nextToken() // point to nessage keyword
			if innerMsg, err = p.parseMessage(recurseDepth + 1); err == nil {
				msg.Messages = append(msg.Messages, innerMsg)
			}
		case token.KindIdentifier:
			var field ast.Field

			if field, err = p.parseField(); err == nil {
				msg.Fields = append(msg.Fields, field)
			}
		default:
			err = gotUnexpected(peek, token.KindOption, token.KindReserved, token.KindIdentifier, token.KindRightBrace)
		}

		if err != nil {
			// TODO report error
			// TODO p.advanceTo(exprEnd)
			return ast.Message{}, err
		}
	}

	if peek.Kind != token.KindRightBrace {
		return ast.Message{}, gotUnexpected(peek, token.KindRightBrace)
	}

	last := p.nextToken()
	msg.Name = id
	msg.ID = p.fm.Merge(token.KindMessage, first.ID, last.ID)
	return msg, nil
}
