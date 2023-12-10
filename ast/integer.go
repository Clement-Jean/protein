package ast

import (
	"fmt"

	"github.com/Clement-Jean/protein/token"
)

type Integer struct {
	ID token.UniqueID
}

func (i Integer) expressionNode()       {}
func (i Integer) GetID() token.UniqueID { return i.ID }
func (i Integer) String() string {
	return fmt.Sprintf("{ type: Int, id: %d }", i.ID)
}
