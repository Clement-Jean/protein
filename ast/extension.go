package ast

import (
	"fmt"

	"github.com/Clement-Jean/protein/token"
)

type ExtensionRange struct {
	Options   []Option
	Ranges    []Range
	OptionsID token.UniqueID
	ID        token.UniqueID
}

func (e ExtensionRange) String() string {
	return fmt.Sprintf("{ type: ExtensionRange, id: %d }", e.ID)
}

type Extend struct {
	Options []Option
	Fields  []Field
	Name    Identifier
	ID      token.UniqueID
}

func (e Extend) String() string {
	return fmt.Sprintf("{ type: Extend, id: %d }", e.ID)
}
