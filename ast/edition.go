package ast

import (
	"fmt"

	"github.com/Clement-Jean/protein/token"
)

type Edition struct {
	ID    token.UniqueID
	Value String
}

func (e Edition) String() string {
	return fmt.Sprintf("{ type: edition, id: %d, value: %s }", e.ID, e.Value.String())
}
