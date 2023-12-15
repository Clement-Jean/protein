package ast

import (
	"fmt"

	"github.com/Clement-Jean/protein/token"
)

type FieldType = uint8

// the order of the following consts is important
// do not change that!
const (
	FieldTypeUnknown FieldType = iota
	FieldTypeDouble
	FieldTypeFloat
	FieldTypeInt64
	FieldTypeUint64
	FieldTypeInt32
	FieldTypeFixed64
	FieldTypeFixed32
	FieldTypeBool
	FieldTypeString
	FieldTypeGroup
	FieldTypeMessage
	FieldTypeBytes
	FieldTypeUint32
	FieldTypeEnum
	FieldTypeSfixed32
	FieldTypeSfixed64
	FieldTypeSint32
	FieldTypeSint64
)

type FieldLabel = uint8

// the order of the following consts is important
// do not change that!
const (
	FieldLabelOptional = iota + 1
	FieldLabelRequired
	FieldLabelRepeated
)

type Field struct {
	Options   []Option
	Name      Identifier
	Tag       Integer
	ID        token.UniqueID
	LabelID   token.UniqueID
	TypeID    token.UniqueID
	OptionsID token.UniqueID
	Label     FieldLabel
	Type      FieldType
}

func (f Field) String() string {
	return fmt.Sprintf("{ type: Field, id: %d, name: %s, tag: %d }", f.ID, f.Name, f.Tag)
}

type Message struct {
	Options       []Option
	Fields        []Field
	ReservedTags  []ReservedTags
	ReservedNames []ReservedNames
	Oneofs        []Oneof
	Enums         []Enum
	Messages      []Message
	Name          Identifier
	ID            token.UniqueID
}

func (m Message) String() string {
	return fmt.Sprintf("{ type: Message, id: %d }", m.ID)
}
