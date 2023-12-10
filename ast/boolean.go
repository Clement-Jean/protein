package ast

import (
	"fmt"

	"github.com/Clement-Jean/protein/token"
)

type Boolean struct {
	ID token.UniqueID
}

func (b Boolean) expressionNode()       {}
func (b Boolean) GetID() token.UniqueID { return b.ID }
func (b Boolean) String() string {
	return fmt.Sprintf("{ type: Boolean, id: %d }", b.ID)
}
