package ast

import (
	"fmt"

	"github.com/Clement-Jean/protein/token"
)

type Package struct {
	Value Identifier
	ID    token.UniqueID
}

func (p Package) String() string {
	return fmt.Sprintf("{ type: package, id: %d, value: %s }", p.ID, p.Value)
}
