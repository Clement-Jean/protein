package ast

import (
	"fmt"

	"github.com/Clement-Jean/protein/token"
)

type Range struct {
	ID    token.UniqueID
	Start Integer
	End   Integer
}

func (r Range) expressionNode()       {}
func (r Range) GetID() token.UniqueID { return r.ID }
func (r Range) String() string {
	return fmt.Sprintf("{ type: Range, id: %d, start: %s, end: %s }", r.ID, r.Start.String(), r.End.String())
}
