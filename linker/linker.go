package linker

import (
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
	units []Unit

	depsIDs   map[string]int // names -> id
	depsNames []string       // id -> names
	depId     int            // current dep id
}

func New(units []Unit) *Linker {
	return &Linker{
		units: units,

		depId:     0,
		depsIDs:   make(map[string]int),
		depsNames: make([]string, len(units)),
	}
}

func (l *Linker) Link() []error {
	// TODO: -I option to recognize that corpus/import_a.proto
	//       and import_a.proto are the same files.

	for i := 0; i < len(l.units); i++ {
		l.depsIDs[l.units[i].File] = l.depId
		l.depsNames[l.depId] = l.units[i].File
		l.depId++
	}

	pkgs := make(map[string]string, len(l.units)) // file -> pkg
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
				l.handlePackage(&pkgs, l.units[i], node.TokIdx)
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

	var multiset []string

	for i := 0; i < len(l.units); i++ {
		pkg := pkgs[l.units[i].File]
		s := []string{pkg}

		for _, node := range l.units[i].Tree {
			switch node.Kind {
			case parser.NodeKindMessageClose, parser.NodeKindEnumClose:
				if len(s) > 0 {
					s = s[:len(s)-1]
				}

			case parser.NodeKindMessageDecl:
				l.handleMessage(&multiset, &s, l.units[i], node.TokIdx)
			case parser.NodeKindEnumDecl:
				l.handleEnum(&multiset, &s, l.units[i], node.TokIdx)
			case parser.NodeKindServiceDecl:
				l.handleService(&multiset, s, l.units[i], node.TokIdx)
			}
		}
	}

	unique := multisetSort(multiset)

	for i := 0; i < len(multiset); i++ {
		print("('", multiset[i], "', ", unique[i], ") ")
	}
	println()

	// TODO check already defined
	// TODO check field types (if identifier)
	// TODO check rpc input and output

	return nil
}
