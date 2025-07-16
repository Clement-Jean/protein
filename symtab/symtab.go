package symtab

import (
	"iter"

	"github.com/Clement-Jean/protein/parser"
	"github.com/Clement-Jean/protein/unit"
)

type Symtab struct {
	Decls map[string]Decl
}

func New() *Symtab {
	return &Symtab{
		Decls: make(map[string]Decl),
	}
}

func (s *Symtab) SearchDecl(name string) (Decl, bool) {
	decl, ok := s.Decls[name]
	return decl, ok
}

func (s *Symtab) AddRef(decl *Decl, fieldName string, fieldTag int64, ref Ref) {
	decl.Fields[fieldName] = ref

	if decl.FieldTags != nil {
		decl.FieldTags[fieldTag] = ref
	}
}

func (s *Symtab) All() iter.Seq2[string, Decl] {
	return func(yield func(string, Decl) bool) {
		for k, v := range s.Decls {
			if !yield(k, v) {
				return
			}
		}
	}
}

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
