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

	depGraph := make([][]int, len(l.units))

	for i := 0; i < len(l.units); i++ {
		for _, node := range l.units[i].Tree {
			switch l.units[i].Toks.TokenInfos[node.TokIdx].Kind {
			case lexer.TokenKindImport:
				l.handleImport(&depGraph, l.units[i], node.TokIdx)
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

	// TODO the parser will have to provide Kinds for
	//      statements. This will help detecting when
	//      we have a field, a message, etc... and this
	//      will let us create symbol tables

	return nil
}
