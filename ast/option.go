package ast

import (
	"fmt"

	"github.com/Clement-Jean/protein/token"
)

type Option struct {
	Value Expression
	Name  Identifier
	ID    token.UniqueID
}

func (o Option) String() string {
	return fmt.Sprintf("{ type: option, id: %d, name: %s, value: %s }", o.ID, o.Name, o.Value)
}
