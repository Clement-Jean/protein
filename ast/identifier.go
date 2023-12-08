package ast

import (
	"fmt"

	"github.com/Clement-Jean/protein/token"
)

type Identifier struct {
	Parts []token.UniqueID
	ID    token.UniqueID
}

func (i Identifier) String() string {
	return fmt.Sprintf("{  type: Identifier, id: %d, parts: %s }", i.ID, fmt.Sprint(i.Parts))
}
