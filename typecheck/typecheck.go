package typecheck

import (
	"errors"
	"os"
	"slices"
	"strings"

	"github.com/Clement-Jean/protein/lexer"
	"github.com/Clement-Jean/protein/parser"
	"github.com/Clement-Jean/protein/source"
	"github.com/bits-and-blooms/bitset"
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

func (tc *TypeChecker) getLineColumn(unit *Unit, offset uint32) (line, col int) {
	line = int(unit.Toks.FindLineIndex(offset))
	lineStart := unit.Toks.LineInfos[line].Start
	col = int(offset - lineStart)
	return line + 1, col + 1
}

func (tc *TypeChecker) registerDep(unit *Unit) {
	tc.depsIDs[unit] = tc.depId
	tc.depsNames[tc.depId] = unit
	tc.depId++
}

func (tc *TypeChecker) checkTypesDeclsRefs(decls, refs *typeMultiset, depGraph [][]int) (errs []error) {
	var used bitset.BitSet
	hasWarning := tc.errorLevel <= ErrorLevelWarning

	for i := 0; i < len(decls.names)-1; i++ {
		name := decls.names[i]
		oldI := i

		for i < len(decls.names)-1 && decls.names[i+1] == name {
			i++
		}

		if oldI != i {
			err := &TypeRedefinedError{
				Name: name.Value(),
			}

			for j := oldI; j <= i; j++ {
				if !decls.kinds[j].IsTypeDef() {
					continue
				}

				unit := decls.units[j]
				offset := decls.offsets[j]
				line, col := tc.getLineColumn(unit, offset)
				err.Files = append(err.Files, unit.File)
				err.Lines = append(err.Lines, line)
				err.Cols = append(err.Cols, col)

				if hasWarning {
					used.Set(uint(j)) // if error don't show warnings...
				}
			}

			errs = append(errs, err)
		}
	}

	// these types are only relevant in the context of this function
	type CacheKey struct {
		unit *Unit
		name string
	}
	type CacheVal struct {
		declIdx         int
		lastNameChecked string
		ok              bool
	}

	cache := make(map[CacheKey]CacheVal)
	for i, ref := range refs.names {
		refName := ref.Value()
		refUnit := refs.units[i]
		refKind := refs.kinds[i]

		var (
			declIdx         int
			lastNameChecked string
			ok              bool
		)

		cacheKey := CacheKey{refUnit, refName}
		if val, hasVal := cache[cacheKey]; hasVal {
			declIdx, lastNameChecked, ok = val.declIdx, val.lastNameChecked, val.ok
		} else {
			declIdx, lastNameChecked, ok = checkUpperScopes(decls, refName)
			cache[cacheKey] = CacheVal{declIdx, lastNameChecked, ok}
		}

		if !ok {
			offset := refs.offsets[i]
			line, col := tc.getLineColumn(refUnit, offset)
			closeIdx := strings.LastIndexByte(refName, ']')

			if closeIdx != -1 {
				refName = refName[closeIdx+1:]
			}

			if lastNameChecked != refName && lastNameChecked != ("."+refName) {
				errs = append(errs, &TypeResolvedNotDefinedError{
					File:         refUnit.File,
					Name:         refName,
					ResolvedName: lastNameChecked,
					Line:         line,
					Col:          col,
				})
			} else {
				errs = append(errs, &TypeNotDefinedError{
					File: refUnit.File,
					Name: refName,
					Line: line,
					Col:  col,
				})
			}
		} else {
			declUnit := decls.units[declIdx]
			declKind := decls.kinds[declIdx]

			if declKind.NotType() {
				offset := refs.offsets[i]
				line, col := tc.getLineColumn(refUnit, offset)
				closeIdx := strings.LastIndexByte(refName, ']')

				if closeIdx != -1 {
					refName = refName[closeIdx+1:]
				}

				errs = append(errs, &NotTypeError{
					Name: refName,
					File: refUnit.File,
					Line: line,
					Col:  col,
				})

				if hasWarning {
					used.Set(uint(declIdx)) // if error don't show warnings...
				}
				continue
			} else if declKind != parser.NodeKindMessageDecl && refKind == parser.NodeKindRPCInputOutput {
				offset := refs.offsets[i]
				line, col := tc.getLineColumn(refUnit, offset)
				closeIdx := strings.LastIndexByte(refName, ']')

				if closeIdx != -1 {
					refName = refName[closeIdx+1:]
				}

				errs = append(errs, &NotMessageTypeError{
					Name: refName,
					File: refUnit.File,
					Line: line,
					Col:  col,
				})

				if hasWarning {
					used.Set(uint(declIdx)) // if error don't show warnings...
				}
				continue
			}

			accessible := declUnit == refUnit || // in same file
				slices.Contains(depGraph[tc.depsIDs[refUnit]], tc.depsIDs[declUnit]) // imported

			if !accessible {
				offset := decls.offsets[declIdx]
				line, col := tc.getLineColumn(declUnit, offset)
				closeIdx := strings.LastIndexByte(refName, ']')

				if closeIdx != -1 {
					refName = refName[closeIdx+1:]
				}

				errs = append(errs, &TypeNotImportedError{
					Name:    refName,
					RefFile: refUnit.File,
					DefFile: declUnit.File,
					Line:    line,
					Col:     col,
				})
			} else if hasWarning {
				used.Set(uint(declIdx))
			}
		}
	}

	if hasWarning && used.Count() != uint(len(decls.names)) {
		for i := 0; i < len(decls.names); i++ {
			ok := used.Test(uint(i))

			if !ok {
				name := decls.names[i].Value()
				unit := decls.units[i]
				offset := decls.offsets[i]
				line, col := tc.getLineColumn(unit, offset)
				closeIdx := strings.LastIndexByte(name, ']')

				if closeIdx != -1 {
					name = name[closeIdx+1:]
				}

				errs = append(errs, &TypeUnusedWarning{
					Name: name,
					File: unit.File,
					Line: line,
					Col:  col,
				})
			}
		}
	}
	return errs
}

func (tc *TypeChecker) checkTypes(depGraph [][]int) []error {
	decls := &typeMultiset{}
	refs := &typeMultiset{}

	for _, unit := range tc.units {
		pkg := tc.pkgs[unit]
		st := []string{pkg} // stack keeping track of type nesting

		for _, node := range unit.Tree {
			declsSize := len(decls.names)
			refsSize := len(refs.names)
			tokIdx := node.TokIdx
			kind := node.Kind

			switch kind {
			case parser.NodeKindMessageClose:
				if len(st) > 0 {
					st = st[:len(st)-1]
				}
				continue

			// DEFS
			case parser.NodeKindMessageDecl:
				tc.handleMessage(decls, &st, unit, tokIdx)
			case parser.NodeKindMessageOneOfDecl:
				tc.handleOneof(decls, st, unit, tokIdx)
			case parser.NodeKindEnumDecl:
				tc.handleEnum(decls, st, unit, tokIdx)
			case parser.NodeKindServiceDecl:
				tc.handleService(decls, st, unit, tokIdx)
			// REFS
			case parser.NodeKindMessageFieldDecl:
				tc.handleField(refs, st, unit, tokIdx)
			case parser.NodeKindMapValue:
				tc.handleMapValue(refs, st, unit, tokIdx)
			case parser.NodeKindRPCInputOutput:
				tc.handleRPC(refs, st, unit, tokIdx)
			// OTHER
			default:
				continue
			}

			if declsSize != len(decls.names) || refsSize != len(refs.names) {
				if node.Kind.IsTypeRef() {
					refs.units = append(refs.units, unit)
					refs.kinds = append(refs.kinds, kind)
				} else if node.Kind.IsTypeDef() {
					decls.units = append(decls.units, unit)
					decls.kinds = append(decls.kinds, kind)
				}
			}
		}
	}

	multisetSort(refs)
	multisetSort(decls)
	return tc.checkTypesDeclsRefs(decls, refs, depGraph)
}

func (tc *TypeChecker) Check() []error {
	// TODO: embed WKT to avoid reparsing them

	var (
		errs      []error
		fatalErrs []error
	)

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

	errs = append(errs, tc.checkTypes(depGraph)...)
	return errs
}
