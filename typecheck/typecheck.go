package typecheck

import (
	"errors"
	"os"
	"slices"
	"strings"
	"unique"

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

type TypeChecker struct {
	srcCreator SourceCreator
	fileCheck  FileExistsCheck
	pkgs       map[*Unit]string

	depId     int
	depsIDs   map[*Unit]int
	depsNames map[int]*Unit

	units        []*Unit
	includePaths []string
}

func New(units []*Unit, opts ...TypeCheckerOpt) *TypeChecker {
	tc := &TypeChecker{
		units: units,

		depId:     0,
		depsIDs:   make(map[*Unit]int, len(units)),
		depsNames: make(map[int]*Unit, len(units)),

		includePaths: []string{""},
		srcCreator:   source.NewFromFile,
		fileCheck: func(path string) bool {
			_, err := os.Stat(path)
			return !errors.Is(err, os.ErrNotExist)
		},
	}

	for _, opt := range opts {
		opt(tc)
	}

	return tc
}

func (tc *TypeChecker) getLineColumn(multiset typeMultiset, i int) (line, col uint32) {
	offset := multiset.offsets[i]
	unit := multiset.units[i]
	line = uint32(unit.Toks.FindLineIndex(offset))
	lineStart := unit.Toks.LineInfos[line].Start
	col = offset - lineStart
	return line, col
}

func (tc *TypeChecker) registerDep(unit *Unit) {
	tc.depsIDs[unit] = tc.depId
	tc.depsNames[tc.depId] = unit
	tc.depId++
}

func (tc *TypeChecker) checkTypeUpperScopes(depGraph [][]int, multiset *typeMultiset, name unique.Handle[string], unit *Unit) (idx int, ok bool) {
	parts := strings.Split(name.Value(), ".")
	namePart := parts[len(parts)-1]

	if len(parts) > 1 {
		cmpFn := func(h unique.Handle[string], s string) int {
			if h.Value() < s {
				return -1
			} else if h.Value() > s {
				return 1
			}

			return 0
		}

		idx := 0
		if len(parts[0]) == 0 {
			idx++
		}

		scope := parts[idx : len(parts)-1]

		for len(scope) != 0 {
			scope = scope[:len(scope)-1] // pop

			newName := strings.Join(scope, ".") + "." + namePart
			idx, ok := slices.BinarySearchFunc(multiset.names, newName, cmpFn)

			if ok {
				if multiset.kinds[idx].IsTypeDef() {
					accessible := multiset.units[idx] == unit || // in same file
						slices.Contains(depGraph[tc.depsIDs[unit]], tc.depsIDs[multiset.units[idx]]) // imported

					if accessible {
						return idx, true
					}
				}

				item := multiset.names[idx]
				idx--
				for idx >= 0 && multiset.names[idx] == item {
					if multiset.kinds[idx].IsTypeDef() {
						accessible := multiset.units[idx] == unit || // in same file
							slices.Contains(depGraph[tc.depsIDs[unit]], tc.depsIDs[multiset.units[idx]]) // imported

						if accessible {
							return idx, true
						}
					}
					idx--
				}

				return 0, false
			}
		}
	}

	return 0, false
}

func (tc *TypeChecker) checkTypes(depGraph [][]int) []error {
	var multiset typeMultiset

	for i := 0; i < len(tc.units); i++ {
		pkg := tc.pkgs[tc.units[i]]
		st := []string{pkg} // stack keeping track of type nesting

		for _, node := range tc.units[i].Tree {
			size := len(multiset.names)

			switch node.Kind {
			case parser.NodeKindMessageClose:
				if len(st) > 0 {
					st = st[:len(st)-1]
				}
				continue

			case parser.NodeKindMessageDecl:
				tc.handleMessage(&multiset, &st, tc.units[i], node.TokIdx)
			case parser.NodeKindMessageFieldDecl:
				tc.handleField(&multiset, st, tc.units[i], node.TokIdx)
			case parser.NodeKindMapValue:
				tc.handleMapValue(&multiset, st, tc.units[i], node.TokIdx)
			case parser.NodeKindEnumDecl:
				tc.handleEnum(&multiset, st, tc.units[i], node.TokIdx)
			case parser.NodeKindServiceDecl:
				tc.handleService(&multiset, st, tc.units[i], node.TokIdx)
			case parser.NodeKindRPCInputOutput:
				tc.handleRPC(&multiset, st, tc.units[i], node.TokIdx)
			default:
				continue
			}

			if size != len(multiset.names) { // added element
				multiset.units = append(multiset.units, tc.units[i])
				multiset.kinds = append(multiset.kinds, node.Kind)
			}
		}
	}

	multisetSort(&multiset)
	n := len(multiset.names)

	var errs []error
	infos := make([][2]int, n) // decls, refs // FIX: can probably make this array smaller

	for i := 0; i < n; i++ {
		var declIdx int
		decls := 0
		refs := 0
		name := multiset.names[i]
		unit := multiset.units[i]

		if multiset.kinds[i].IsTypeDef() {
			declIdx = i
			decls++
		}

		for i = i + 1; i < n && multiset.names[i] == name; i++ {
			if multiset.kinds[i].IsTypeRef() {
				accessible := multiset.units[i] == unit || // in same file
					slices.Contains(depGraph[tc.depsIDs[multiset.units[i]]], tc.depsIDs[unit]) // imported

				if accessible {
					refs++
				}
			} else if multiset.kinds[i].IsTypeDef() {
				decls++
				// FIX: maybe we can break here? because we know it redefined
			}
		}
		i--

		if decls == 0 {
			if idx, ok := tc.checkTypeUpperScopes(depGraph, &multiset, name, unit); ok {
				// found in upper scopes
				infos[idx][1]++
				continue
			}

			line, col := tc.getLineColumn(multiset, i)
			println("file", multiset.units[i].File, "col", col, "line", line+1, "name", name.Value())

			// TODO: improve error message to show the line of the ref
			errs = append(errs, &TypeNotDefinedError{Name: name.Value()})
		} else {
			infos[declIdx][0] += decls
			infos[declIdx][1] += refs

			if decls > 1 {
				// backtracks to find the multiple defs
				for j := i; j >= 0 && multiset.names[j] == name; j-- {
					if !multiset.kinds[j].IsTypeDef() {
						continue
					}

					//line, col := tc.getLineColumn(multiset, j)
					//println("file", multiset.units[i].File, "col", col, "line", line+1, "name", name.Value())
				}

				// TODO: improve error message to show the 2+ declarations
				errs = append(errs, &TypeRedefinedError{Name: name.Value()})
			}
		}

		// println("type", name.Value(), "has", decls, "decls and", refs, "refs")
	}

	// TODO: skip the following loop if the error level is greater than WARNING

	// we can only check the uses of a type here because
	// checkTypeUpperScopes could use a type later/earlier
	// than when we check its definition, it could lead to
	// incorrect warnings.
	for i := 0; i < n; i++ {
		if multiset.kinds[i].IsTypeDef() && infos[i][0] == 1 && infos[i][1] == 0 {
			name := multiset.names[i].Value()
			//line, col := tc.getLineColumn(multiset, i)
			//println("file", multiset.units[i].File, "col", col, "line", line+1, "name", name)

			// TODO: improve error message to show the line of the ref
			errs = append(errs, &TypeUnusedWarning{Name: name})
			//continue
		}

		print("(\"", multiset.names[i].Value(), "\", ", multiset.kinds[i].String(), "),")
	}
	println()

	// TODO check rpc input and output
	return errs
}

func (tc *TypeChecker) Check() []error {
	// TODO: embed WKT to avoid reparsing them

	var errs []error

	for j := 0; j < len(tc.units); j++ {
		tc.registerDep(tc.units[j])
	}

	tc.pkgs = make(map[*Unit]string, len(tc.units))
	unitsLen := len(tc.units)
	depGraph := make([][]int, unitsLen)
	i := 0

	for true {
		// unfortunately imports can be placed anywhere.
		// This means we need to resolve all of them first
		// before being able to resolve types.
		for j := i; j < len(tc.units); j++ {
			for _, node := range tc.units[j].Tree {
				switch node.Kind {
				case parser.NodeKindImportStmt:
					if err := tc.handleImport(&depGraph, tc.units[j], node.TokIdx); err != nil {
						errs = append(errs, err)
					}
				}
			}
		}

		if len(errs) != 0 {
			return errs
		}

		if len(tc.units) == unitsLen { // all imports handled
			break
		}

		errs = append(errs, tc.handleUnknownImports(i)...)
		if len(errs) != 0 {
			return errs
		}

		i = unitsLen // recheck only newly added
		unitsLen = len(tc.units)
	}

	// unfortunately packages can be placed anywhere.
	// This means we need to resolve all of them first
	// before being able to resolve types.
	for j, unit := range tc.units {
		for _, node := range unit.Tree {
			switch node.Kind {
			case parser.NodeKindPackageStmt:
				if _, ok := tc.pkgs[tc.units[j]]; ok {
					errs = append(errs, &PackageMultipleDefError{File: tc.units[j].File})
					break
				}

				tc.handlePackage(tc.pkgs, tc.units[j], node.TokIdx)
			}
		}
	}

	errs = append(errs, tc.checkImportCycles(depGraph)...)
	if len(errs) != 0 {
		return errs
	}

	errs = append(errs, tc.checkTypes(depGraph)...)
	return errs
}
