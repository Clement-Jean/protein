package ast

import (
	"fmt"

	"github.com/Clement-Jean/protein/token"
)

type Oneof struct {
	Options []Option
	Fields  []Field
	Name    Identifier
	ID      token.UniqueID
}

func (o Oneof) String() string {
	return fmt.Sprintf("{ type: Oneof, id: %d }", o.ID)
}
