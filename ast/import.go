package ast

import (
	"fmt"

	"github.com/Clement-Jean/protein/token"
)

type Import struct {
	ID       token.UniqueID
	Value    String
	IsPublic bool
	IsWeak   bool
}

func (i Import) String() string {
	return fmt.Sprintf("{ type: import, id: %d, value: %s }", i.ID, i.Value.String())
}
