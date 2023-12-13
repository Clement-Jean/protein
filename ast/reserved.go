package ast

import (
	"fmt"

	"github.com/Clement-Jean/protein/token"
)

type ReservedTags struct {
	Items []Range
	ID    token.UniqueID
}

func (r ReservedTags) GetID() token.UniqueID { return r.ID }
func (r ReservedTags) String() string {
	return fmt.Sprintf("{ type: ReservedTags, id: %d }", r.ID)
}

type ReservedNames struct {
	Items []String
	ID    token.UniqueID
}

func (r ReservedNames) GetSpan() token.UniqueID { return r.ID }
func (r ReservedNames) String() string {
	return fmt.Sprintf("{ type: ReservedNames, id: %d }", r.ID)
}
