package ast

import (
	"fmt"

	"github.com/Clement-Jean/protein/token"
)

type EnumValue struct {
	Options   []Option
	Name      Identifier
	ID        token.UniqueID
	Tag       Integer
	OptionsID token.UniqueID
}

func (ev EnumValue) GetID() token.UniqueID { return ev.ID }
func (ev EnumValue) String() string {
	return fmt.Sprintf("{ type: enum_value, id: %d, name: %s, tag: %d }", ev.ID, ev.Name, ev.Tag)
}

type Enum struct {
	Options       []Option
	Values        []EnumValue
	ReservedTags  []ReservedTags
	ReservedNames []ReservedNames
	Name          Identifier
	ID            token.UniqueID
}

func (e Enum) GetID() token.UniqueID { return e.ID }
func (e Enum) String() string {
	return fmt.Sprintf("{ type: enum, id: %d }", e.ID)
}
