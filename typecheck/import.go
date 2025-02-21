package typecheck

import (
	"path/filepath"
	"slices"
	"strings"

	"github.com/Clement-Jean/protein/lexer"
	"github.com/Clement-Jean/protein/parser"
)

func (tc *TypeChecker) handleImport(depGraph *[][]int, unit *Unit, idx uint32) error {
	if unit.Toks.TokenInfos[idx+1].Kind == lexer.TokenKindIdentifier {
		idx += 1
	}
	idx += 1

	start := unit.Toks.TokenInfos[idx].Offset
	end := unit.Toks.TokenInfos[idx+1].Offset
	file := string(unit.Buffer.Range(start, end))
	file = strings.Trim(file, " \"'")

	var to *Unit

	for i := 0; i < len(tc.units); i++ {
		for j := 0; j < len(tc.includePaths); j++ {
			path := file

			if !strings.HasPrefix(file, tc.includePaths[j]) {
				path = filepath.Join(tc.includePaths[j], file)
			}

			if ok := tc.fileCheck(path); !ok {
				continue
			}

			if tc.units[i].File == path {
				to = &tc.units[i]
				goto found
			}
		}
	}

found:

	if to == nil {
		// add import to be parsed late (see: handleUnknownImports)
		tc.units = append(tc.units, Unit{File: file})
		to = &tc.units[len(tc.units)-1]
		tc.registerDep(to)
		*depGraph = append(*depGraph, make([]int, 0))
	}

	from := tc.depsIDs[unit]
	(*depGraph)[from] = append((*depGraph)[from], tc.depsIDs[to])
	return nil
}

// parse the file included but not added as inputs
func (tc *TypeChecker) handleUnknownImports(offset int) []error {
	var errs []error

	for i := offset; i < len(tc.units); i++ {
		if tc.units[i].Tree != nil { // already handled
			continue
		}

		var (
			err error          = nil
			lex *lexer.Lexer   = nil
			p   *parser.Parser = nil
		)

		for _, includePath := range tc.includePaths {
			path := tc.units[i].File

			if !strings.HasPrefix(tc.units[i].File, includePath) {
				path = filepath.Join(includePath, tc.units[i].File)
			}

			if ok := tc.fileCheck(path); !ok {
				continue // try other include paths
			}

			tc.units[i].Buffer, err = tc.srcCreator(path)
			if err != nil {
				panic(err) // TODO better handling
			}

			lex, err = lexer.NewFromSource(tc.units[i].Buffer)
			if err != nil {
				panic(err) // TODO better handling
			}

			tc.units[i].Toks, errs = lex.Lex()
			if len(errs) != 0 {
				panic(errs) // TODO better handling
			}

			p = parser.New(tc.units[i].Toks)
			tc.units[i].Tree, errs = p.Parse()
			if len(errs) != 0 {
				panic(errs) // TODO better handling
			}

			break
		}

		if tc.units[i].Tree == nil {
			errs = append(errs, &ImportFileNotFoundError{
				File:         tc.units[i].File,
				IncludePaths: tc.includePaths,
			})
		}
	}

	return errs
}

func (tc *TypeChecker) importCycle(depGraph [][]int) []int {
	n := len(depGraph)
	colors := make([]uint8, n)
	parent := make([]int, n)
	cycleStart := -1
	cycleEnd := -1

	for i := 0; i < n; i++ {
		parent[i] = -1
	}

	dfs := func(v int) bool {
		var s []int

		s = append(s, v) // push

		for len(s) != 0 {
			last := len(s) - 1
			v = s[last]

			if colors[v] != 1 {
				colors[v] = 1 // GREY

				for _, w := range depGraph[v] {
					switch colors[w] {
					case 0: // WHITE
						parent[w] = v
						s = append(s, w) // push
					case 1: // GREY
						cycleStart = w
						cycleEnd = v
						return true
					}
				}
			} else if colors[v] == 1 {
				s = s[:last]  // pop
				colors[v] = 2 // BLACK
			}
		}

		return false
	}

	for i := 0; i < n; i++ {
		if colors[i] == 0 && dfs(i) { // WHITE
			break
		}
	}

	if cycleStart != -1 {
		cycle := []int{cycleStart}
		for v := cycleEnd; v != cycleStart; v = parent[v] {
			cycle = append(cycle, v)
		}
		cycle = append(cycle, cycleStart)
		slices.Reverse(cycle)
		return cycle
	}

	return nil
}

func (tc *TypeChecker) checkImportCycles(depGraph [][]int) []error {
	if cycle := tc.importCycle(depGraph); len(cycle) != 0 {
		var err ImportCycleError
		for _, v := range cycle[:len(cycle)-1] {
			err.Files = append(err.Files, tc.depsNames[v].File)
		}
		return []error{&err}
	}

	return nil
}
