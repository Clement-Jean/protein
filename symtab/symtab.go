package symtab

import (
	"github.com/Clement-Jean/protein/parser"
	"github.com/Clement-Jean/protein/unit"
)

type Symtab = map[string]Decl

type Decl struct {
	Unit      *unit.Unit
	Line, Col int
	Type      parser.NodeKind
	Name      string
	Fields    map[string]Ref
	FieldTags map[int64]Ref // FIX change to uint64
}

type Ref struct {
	// FIX having both Type and TypeName is kind of redundant
	//     could we somehow merge them?
	Unit      *unit.Unit
	Line, Col int
	Type      parser.NodeKind
	TypeName  string
	Tag       int64 // FIX change to uint64
}

// TODO to descriptor
// TODO from descriptor
