package ast

import (
	"fmt"

	"github.com/Clement-Jean/protein/token"
)

type String struct {
	ID token.UniqueID
}

func (s String) String() string {
	return fmt.Sprintf("{ type: String, id: %d }", s.ID)
}
