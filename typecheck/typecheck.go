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

	errorLevel ErrorLevel
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

		errorLevel: ErrorLevelUndefined,
	}

	for _, opt := range opts {
		opt(tc)
	}

	return tc
}

func (tc *TypeChecker) getLineColumn(unit *Unit, offset uint32) (line, col uint32) {
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

			unit := multiset.units[i]
			offset := multiset.offsets[i]
			line, col := tc.getLineColumn(unit, offset)
			errs = append(errs, &TypeNotDefinedError{
				Name: name.Value(),
				File: multiset.units[i].File,
				Line: int(line + 1),
				Col:  int(col + 1),
			})
		} else {
			infos[declIdx][0] += decls
			infos[declIdx][1] += refs

			if decls > 1 {
				// backtracks to find the multiple defs
				var (
					files []string
					lines []int
					cols  []int
				)
				for j := i; j >= 0 && multiset.names[j] == name; j-- {
					if !multiset.kinds[j].IsTypeDef() {
						continue
					}

					unit := multiset.units[j]
					offset := multiset.offsets[j]
					line, col := tc.getLineColumn(unit, offset)
					files = append(files, multiset.units[j].File)
					lines = append(lines, int(line+1))
					cols = append(cols, int(col+1))
				}

				errs = append(errs, &TypeRedefinedError{
					Name:  name.Value(),
					Files: files,
					Lines: lines,
					Cols:  cols,
				})
			}
		}
	}

	if tc.errorLevel <= ErrorLevelWarning {
		// we can only check the uses of a type here because
		// checkTypeUpperScopes could use a type later/earlier
		// than when we check its definition, it could lead to
		// incorrect warnings.
		for i := 0; i < n; i++ {
			decls := infos[i][0]
			refs := infos[i][1]

			if multiset.kinds[i].IsTypeDef() && decls == 1 && refs == 0 {
				name := multiset.names[i].Value()
				unit := multiset.units[i]
				offset := multiset.offsets[i]
				line, col := tc.getLineColumn(unit, offset)
				errs = append(errs, &TypeUnusedWarning{
					Name: name,
					File: multiset.units[i].File,
					Line: int(line + 1),
					Col:  int(col + 1),
				})
				continue
			}

			print("(\"", multiset.names[i].Value(), "\", ", multiset.kinds[i].String(), "),")
		}
		println()
	}

	return errs
}

func (tc *TypeChecker) Check() []error {
	// TODO: embed WKT to avoid reparsing them

	var errs []error
	var fatalErrs []error

	for j := 0; j < len(tc.units); j++ {
		tc.registerDep(tc.units[j])
	}

	tc.pkgs = make(map[*Unit]string, len(tc.units))
	unitsLen := len(tc.units)
	depGraph := make([][]int, unitsLen)
	i := 0

	for true {
		// unfortunately imports and packages can be placed anywhere.
		// This means we need to resolve all of them first
		// before being able to resolve types.
		for j := i; j < len(tc.units); j++ {
			for _, node := range tc.units[j].Tree {
				switch node.Kind {
				case parser.NodeKindImportStmt:
					if err := tc.handleImport(&depGraph, tc.units[j], node.TokIdx); err != nil {
						errs = append(errs, err)
					}
				case parser.NodeKindPackageStmt:
					if _, ok := tc.pkgs[tc.units[j]]; ok {
						fatalErrs = append(fatalErrs, &PackageMultipleDefError{File: tc.units[j].File})
						break
					}

					tc.handlePackage(tc.pkgs, tc.units[j], node.TokIdx)
				}
			}
		}

		if len(tc.units) == unitsLen { // all imports handled
			break
		}

		errs = append(errs, tc.handleUnknownImports(i)...)
		i = unitsLen // recheck only newly added
		unitsLen = len(tc.units)
	}

	fatalErrs = append(fatalErrs, tc.checkImportCycles(depGraph)...)
	if len(fatalErrs) != 0 {
		return slices.Concat(fatalErrs, errs)
	}

	// TODO "a.b.c.D" seems to be defined in "test_d.proto", which is not imported by "test_a.proto".  To use it here, please add the necessary import.

	errs = append(errs, tc.checkTypes(depGraph)...)
	return errs
}
