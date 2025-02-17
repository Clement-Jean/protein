package linker

import (
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/Clement-Jean/protein/lexer"
	"github.com/Clement-Jean/protein/parser"
	"github.com/Clement-Jean/protein/source"
)

func (l *Linker) handleImport(depGraph *[][]int, unit Unit, idx uint32) {
	if unit.Toks.TokenInfos[idx+1].Kind == lexer.TokenKindIdentifier {
		idx += 1
	}
	idx += 1

	start := unit.Toks.TokenInfos[idx].Offset
	end := unit.Toks.TokenInfos[idx+1].Offset
	file := string(unit.Buffer.Range(start, end))
	file = strings.Trim(file, " \"'")

	if _, ok := l.depsIDs[file]; !ok {
		l.depsIDs[file] = l.depId
		l.depId++

		var (
			err  error                  = nil
			errs []error                = nil
			s    *source.Buffer         = nil
			lex  *lexer.Lexer           = nil
			tb   *lexer.TokenizedBuffer = nil
			p    *parser.Parser         = nil
			pt   parser.ParseTree       = nil
		)

		for i := 0; pt == nil && i < len(l.includePaths); i++ {
			path := filepath.Join(l.includePaths[i], file)

			if _, err = os.Stat(path); errors.Is(err, os.ErrNotExist) {
				continue
			}

			s, err = source.NewFromFile(path)
			if err != nil {
				panic(err) // TODO better handling
			}

			lex, err = lexer.NewFromSource(s)
			if err != nil {
				panic(err) // TODO better handling
			}

			tb, errs = lex.Lex()
			if len(errs) != 0 {
				panic(errs) // TODO better handling
			}

			p = parser.New(tb)
			pt, errs = p.Parse()
			if len(errs) != 0 {
				panic(errs) // TODO better handling
			}
		}

		if pt == nil {
			panic(file + " not found") // TODO better handling
		}

		l.units = append(l.units, Unit{
			File:   file,
			Buffer: s,
			Toks:   tb,
			Tree:   pt,
		})
		*depGraph = append(*depGraph, make([]int, 0))
		l.depsNames = append(l.depsNames, file)
	}

	from := l.depsIDs[unit.File]
	(*depGraph)[from] = append((*depGraph)[from], l.depsIDs[file])
}

func (l *Linker) importCycle(depGraph [][]int) []int {
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
		if colors[i] == 0 && dfs(i) {
			break
		}
	}

	var cycle []int
	if cycleStart != -1 {
		cycle = append(cycle, cycleStart)
		for v := cycleEnd; v != cycleStart; v = parent[v] {
			cycle = append(cycle, v)
		}
		cycle = append(cycle, cycleStart)
		slices.Reverse(cycle)
	}

	return cycle
}
