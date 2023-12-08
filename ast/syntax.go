package ast

import (
	"fmt"

	"github.com/Clement-Jean/protein/token"
)

type Syntax struct {
	ID    token.UniqueID
	Value String
}

func (s Syntax) String() string {
	return fmt.Sprintf("{ type: syntax, id: %d, value: %s }", s.ID, s.Value.String())
}
