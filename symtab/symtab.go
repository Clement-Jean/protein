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
	Fields    map[string]Ref
	FieldTags map[int64]Ref
	Name      string
	Line      int
	Col       int
	Type      parser.NodeKind
}

type Ref struct {
	Unit     *unit.Unit
	TypeName string
	Line     int
	Col      int
	Tag      int64
	Type     parser.NodeKind
}

// TODO to descriptor
// TODO from descriptor
