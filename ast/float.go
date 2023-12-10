package ast

import (
	"fmt"

	"github.com/Clement-Jean/protein/token"
)

type Float struct {
	ID token.UniqueID
}

func (f Float) expressionNode()       {}
func (f Float) GetID() token.UniqueID { return f.ID }
func (f Float) String() string {
	return fmt.Sprintf("{  type: Float, id: %d }", f.ID)
}
