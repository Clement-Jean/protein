package typecheck

import (
	"errors"
	"os"
	"slices"
	"strings"

	"github.com/Clement-Jean/protein/parser"
	"github.com/Clement-Jean/protein/source"
	"github.com/Clement-Jean/protein/symtab"
	"github.com/Clement-Jean/protein/unit"
)

type TypeChecker struct {
	srcCreator SourceCreator
	fileCheck  FileExistsCheck
	pkgs       map[*unit.Unit]string

	depId     int
	depsIDs   map[*unit.Unit]int
	depsNames map[int]*unit.Unit

	units        []*unit.Unit
	includePaths []string

	errorLevel ErrorLevel
}

func New(units []*unit.Unit, opts ...TypeCheckerOpt) *TypeChecker {
	tc := &TypeChecker{
		units: units,

		depId:     0,
		depsIDs:   make(map[*unit.Unit]int, len(units)),
		depsNames: make(map[int]*unit.Unit, len(units)),

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

func (tc *TypeChecker) getLineColumn(unit *unit.Unit, offset uint32) (line, col int) {
	line = int(unit.Toks.FindLineIndex(offset))
	lineStart := unit.Toks.LineInfos[line].Start
	col = int(offset - lineStart)
	return line + 1, col + 1
}

func (tc *TypeChecker) registerDep(unit *unit.Unit) {
	tc.depsIDs[unit] = tc.depId
	tc.depsNames[tc.depId] = unit
	tc.depId++
}

// TODO redefined error
// TODO RPC "is not message type"
func (tc *TypeChecker) checkTypesDeclsRefs(sym symtab.Symtab, depGraph [][]int) (errs []error) {
	// these types are only relevant in the context of this function
	type CacheKey struct {
		unit *unit.Unit
		name string
	}
	type CacheVal struct {
		decl            symtab.Decl
		lastNameChecked string
		ok              bool
	}

	cache := make(map[CacheKey]CacheVal)

	for _, symbol := range sym {
		for _, ref := range symbol.Fields {
			if ref.Type != parser.NodeKindUndefined || !strings.HasPrefix(ref.TypeName, ".") {
				continue
			}

			var (
				lastNameChecked string
				decl            symtab.Decl
				ok              bool
			)

			cacheKey := CacheKey{symbol.Unit, ref.TypeName}
			if val, hasVal := cache[cacheKey]; hasVal {
				lastNameChecked, ok = val.lastNameChecked, val.ok
			} else {
				lastNameChecked, decl, ok = checkUpperScopes(sym, ref.TypeName)
				cache[cacheKey] = CacheVal{decl, lastNameChecked, ok}
			}

			if !ok { // not found
				closeIdx := strings.LastIndexByte(ref.TypeName, ']')

				if closeIdx != -1 {
					ref.TypeName = ref.TypeName[closeIdx+1:]
				}

				if lastNameChecked != ref.TypeName && lastNameChecked != ("."+ref.TypeName) {
					errs = append(errs, &TypeResolvedNotDefinedError{
						File:         symbol.Unit.File,
						Name:         ref.TypeName,
						ResolvedName: lastNameChecked,
						Line:         ref.Line,
						Col:          ref.Col,
					})
				} else {
					errs = append(errs, &TypeNotDefinedError{
						File: symbol.Unit.File,
						Name: ref.TypeName,
						Line: ref.Line,
						Col:  ref.Col,
					})
				}
			} else {
				if decl.Type.NotType() {
					closeIdx := strings.LastIndexByte(ref.TypeName, ']')

					if closeIdx != -1 {
						ref.TypeName = ref.TypeName[closeIdx+1:]
					}

					errs = append(errs, &NotTypeError{
						File: symbol.Unit.File,
						Name: ref.TypeName,
						Line: ref.Line,
						Col:  ref.Col,
					})
					continue
				} else if decl.Type != parser.NodeKindMessageDecl && ref.Type == parser.NodeKindRPCInputOutput {
					closeIdx := strings.LastIndexByte(ref.TypeName, ']')

					if closeIdx != -1 {
						ref.TypeName = ref.TypeName[closeIdx+1:]
					}

					errs = append(errs, &NotMessageTypeError{
						File: symbol.Unit.File,
						Name: ref.TypeName,
						Line: ref.Line,
						Col:  ref.Col,
					})
					continue
				}

				accessible := decl.Unit == ref.Unit || // in same file
					slices.Contains(depGraph[tc.depsIDs[ref.Unit]], tc.depsIDs[decl.Unit]) // imported

				if !accessible {
					closeIdx := strings.LastIndexByte(ref.TypeName, ']')

					if closeIdx != -1 {
						ref.TypeName = ref.TypeName[closeIdx+1:]
					}

					errs = append(errs, &TypeNotImportedError{
						Name:    ref.TypeName,
						DefFile: decl.Unit.File,
						RefFile: symbol.Unit.File,
						Line:    ref.Line,
						Col:     ref.Col,
					})
				}
			}
		}
	}
	return errs
}

func (tc *TypeChecker) checkTypes(depGraph [][]int) (symtab.Symtab, []error) {
	sym := make(map[string]symtab.Decl)
	var errs []error

	for _, unit := range tc.units {
		pkg := tc.pkgs[unit]
		var st []string // stack keeping track of type nesting

		if len(pkg) != 0 {
			st = append(st, pkg)
		}

		for i := 0; i < len(unit.Tree); i++ {
			node := unit.Tree[i]
			tokIdx := node.TokIdx
			kind := node.Kind

			switch kind {
			case parser.NodeKindMessageClose,
				parser.NodeKindEnumClose,
				parser.NodeKindServiceClose:
				if len(st) > 0 {
					st = st[:len(st)-1]
				}
				continue

			// DEFS
			case parser.NodeKindMessageDecl:
				if err := tc.handleMessage(sym, &st, unit, tokIdx); err != nil {
					for i < len(unit.Tree) && unit.Tree[i].Kind != parser.NodeKindMessageClose {
						i++
					}
					errs = append(errs, err)
				}
			case parser.NodeKindOneOfDecl:
				if err := tc.handleOneof(sym, st, unit, tokIdx); err != nil {
					for i < len(unit.Tree) && unit.Tree[i].Kind != parser.NodeKindOneofClose {
						i++
					}
					errs = append(errs, err)
				}
			case parser.NodeKindMapDecl:
				if err := tc.handleMapDecl(sym, st, unit, tokIdx); err != nil {
					errs = append(errs, err)
				}
			case parser.NodeKindEnumDecl:
				if err := tc.handleEnum(sym, &st, unit, tokIdx); err != nil {
					for i < len(unit.Tree) && unit.Tree[i].Kind != parser.NodeKindEnumClose {
						i++
					}
					errs = append(errs, err)
				}
			case parser.NodeKindServiceDecl:
				if err := tc.handleService(sym, &st, unit, tokIdx); err != nil {
					for i < len(unit.Tree) && unit.Tree[i].Kind != parser.NodeKindServiceClose {
						i++
					}
					errs = append(errs, err)
				}
			case parser.NodeKindRPCDecl:
				if err := tc.handleRPCDecl(sym, st, unit, tokIdx); err != nil {
					errs = append(errs, err)
				}

			// REFS
			case parser.NodeKindMessageFieldDecl:
				if err := tc.handleField(sym, st, unit, tokIdx); err != nil {
					errs = append(errs, err)
				}
			case parser.NodeKindMapValue:
				if err := tc.handleMapValue(sym, st, unit, tokIdx); err != nil {
					errs = append(errs, err)
				}
			case parser.NodeKindRPCInputOutput:
				if err := tc.handleRPCInputOutput(sym, st, unit, tokIdx); err != nil {
					errs = append(errs, err)
				}
			case parser.NodeKindEnumValueDecl:
				if err := tc.handleValue(sym, st, unit, tokIdx); err != nil {
					errs = append(errs, err)
				}
			// OTHER
			default:
				continue
			}
		}
	}

	for _, value := range sym {
		switch value.Type {
		case parser.NodeKindEnumDecl:
			if len(value.Fields) == 0 {
				errs = append(errs, &EnumCannotBeEmptyError{
					File: value.Unit.File,
					Line: value.Line,
					Col:  value.Col,
					Name: value.Name,
				})
			}
		case parser.NodeKindOneOfDecl:
			// TODO
			// if len(value.Fields) < 2 {
			// 	errs = append(errs, &OneofCannotHaveLessThan2FieldsError{
			// 		File: value.Unit.File,
			// 		Line: value.Line,
			// 		Col:  value.Col,
			// 		Name: value.Name,
			// 	})
			// }
		}
	}

	errs = append(errs, tc.checkTypesDeclsRefs(sym, depGraph)...)
	return sym, errs
}

func (tc *TypeChecker) Check() (symtab.Symtab, []error) {
	// TODO: embed WKT to avoid reparsing them

	var (
		errs      []error
		fatalErrs []error
	)

	for j := range len(tc.units) {
		tc.registerDep(tc.units[j])
	}

	tc.pkgs = make(map[*unit.Unit]string, len(tc.units))
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
						errs = append(errs, err...)
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
		return nil, slices.Concat(fatalErrs, errs)
	}

	sym, err := tc.checkTypes(depGraph)
	errs = append(errs, err...)
	return sym, errs
}
