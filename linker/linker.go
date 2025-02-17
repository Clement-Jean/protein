package linker

import (
	"fmt"

	"github.com/Clement-Jean/protein/lexer"
	"github.com/Clement-Jean/protein/parser"
	"github.com/Clement-Jean/protein/source"
)

type Unit struct {
	File   string
	Buffer *source.Buffer
	Toks   *lexer.TokenizedBuffer
	Tree   parser.ParseTree
}

type Linker struct {
	units        []Unit
	includePaths []string

	depsIDs   map[string]int // names -> id
	depsNames []string       // id -> names
	depId     int            // current dep id

	pkgs map[string]string // file -> pkg
}

func New(units []Unit, opts ...LinkerOpt) *Linker {
	l := &Linker{
		units: units,

		depId:     0,
		depsIDs:   make(map[string]int),
		depsNames: make([]string, len(units)),
	}

	for _, opt := range opts {
		opt(l)
	}

	return l
}

func (l *Linker) checkTypes() []error {
	var multiset typeMultiset

	for i := 0; i < len(l.units); i++ {
		pkg := l.pkgs[l.units[i].File]
		st := []string{pkg} // stack keeping track of type nesting

		for _, node := range l.units[i].Tree {
			size := len(multiset.names)

			switch node.Kind {
			case parser.NodeKindMessageClose:
				if len(st) > 0 {
					st = st[:len(st)-1]
				}

			case parser.NodeKindMessageDecl:
				l.handleMessage(&multiset, &st, l.units[i], node.TokIdx)
			case parser.NodeKindMessageFieldDecl:
				l.handleField(&multiset, st, l.units[i], node.TokIdx)
			case parser.NodeKindMapValue:
				l.handleMapValue(&multiset, st, l.units[i], node.TokIdx)
			case parser.NodeKindEnumDecl:
				l.handleEnum(&multiset, st, l.units[i], node.TokIdx)
			case parser.NodeKindServiceDecl:
				l.handleService(&multiset, st, l.units[i], node.TokIdx)
			}

			if size != len(multiset.names) { // added element
				multiset.kinds = append(multiset.kinds, node.Kind)
			}
		}
	}

	multisetSort(multiset)
	n := len(multiset.names)

	var errs []error

	for i := 0; i < n; i++ {
		decls := 0
		refs := 0
		name := multiset.names[i]

		if multiset.kinds[i].IsTypeDef() {
			decls++
		} else if multiset.kinds[i].IsTypeRef() {
			refs++
		}

		for i = i + 1; i < n && multiset.names[i] == name; i++ {
			if multiset.kinds[i].IsTypeDef() {
				decls++
			} else if multiset.kinds[i].IsTypeRef() {
				refs++
			}
		}
		i--

		if decls == 0 {
			// TODO: improve error message to show the line of the ref
			errs = append(errs, fmt.Errorf("%s is not defined", name.Value()))
		} else if refs == 0 {
			// TODO: make this a warning
			errs = append(errs, fmt.Errorf("%s is not used", name.Value()))
		} else if decls > 2 {
			// TODO: improve error message to show the 2+ declarations
			errs = append(errs, fmt.Errorf("%s is redefined", name.Value()))
		}

		println("type", name.Value(), "has", decls, "decls and", refs, "refs")
	}

	for i := 0; i < n; i++ {
		print("(\"", multiset.names[i].Value(), "\"),")
	}
	println()

	// TODO check already defined
	// TODO check field types (if identifier)
	// TODO check rpc input and output
	return errs
}

func (l *Linker) checkImportCycles() []error {
	for i := 0; i < len(l.units); i++ {
		l.depsIDs[l.units[i].File] = l.depId
		l.depsNames[l.depId] = l.units[i].File
		l.depId++
	}

	l.pkgs = make(map[string]string, len(l.units))
	depGraph := make([][]int, len(l.units))

	// unfortunately imports and packages can be placed
	// anywhere. This means we need to resolve all of
	// them first before being able to resolve types.
	for i := 0; i < len(l.units); i++ {
		for _, node := range l.units[i].Tree {
			switch node.Kind {
			case parser.NodeKindImportStmt:
				l.handleImport(&depGraph, l.units[i], node.TokIdx)
			case parser.NodeKindPackageStmt:
				l.handlePackage(&l.pkgs, l.units[i], node.TokIdx)
			}
		}
	}

	if cycle := l.importCycle(depGraph); len(cycle) != 0 {
		var files []string
		for _, v := range cycle[:len(cycle)-1] {
			files = append(files, l.depsNames[v])
		}

		return []error{&ImportCycleError{Files: files}}
	}

	return nil
}

func (l *Linker) Link() []error {
	// TODO: -I option to recognize that corpus/import_a.proto
	//       and import_a.proto are the same files.

	var errs []error
	errs = append(errs, l.checkImportCycles()...)
	errs = append(errs, l.checkTypes()...)
	return errs
}
