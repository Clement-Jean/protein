package typecheck

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/Clement-Jean/protein/parser"
	"github.com/Clement-Jean/protein/symtab"
	"github.com/Clement-Jean/protein/unit"
)

func (tc *TypeChecker) handleEnum(sym *symtab.Symtab, scope *[]string, unit *unit.Unit, idx uint32) error {
	start := unit.Toks.TokenInfos[idx].Offset
	end := unit.Toks.TokenInfos[idx+1].Offset
	line, col := tc.getLineColumn(unit, start)
	name := strings.TrimSpace(string(unit.Buffer.Range(start, end)))
	prefix := strings.Join(*scope, ".")

	if len(prefix) != 0 && !strings.HasPrefix(prefix, ".") {
		prefix = "." + prefix
	}

	(*scope) = append((*scope), name)

	fullName := fmt.Sprintf("%s.%s", prefix, name)

	if decl, ok := sym.SearchDecl(fullName); ok {
		return &TypeRedefinedError{
			Name:  fullName,
			Files: []string{unit.File, decl.Unit.File},
			Lines: []int{line, decl.Line},
			Cols:  []int{col, decl.Col},
		}
	}

	sym.Decls[fullName] = symtab.Decl{
		Unit:      unit,
		Line:      line,
		Col:       col,
		Type:      parser.NodeKindEnumDecl,
		Name:      name,
		Fields:    make(map[string]symtab.Ref),
		FieldTags: make(map[int64]symtab.Ref),
	}
	return nil
}

func (tc *TypeChecker) handleValue(sym *symtab.Symtab, scope []string, unit *unit.Unit, idx uint32) error {
	start := unit.Toks.TokenInfos[idx].Offset
	end := unit.Toks.TokenInfos[idx+1].Offset
	line, col := tc.getLineColumn(unit, start)
	valueName := strings.TrimSpace(string(unit.Buffer.Range(start, end)))

	endOffset := unit.Toks.TokenInfos[idx+2].Offset
	endNextOffset := unit.Toks.TokenInfos[idx+3].Offset
	valueTag, _ := strconv.ParseInt(string(unit.Buffer.Range(endOffset, endNextOffset)), 10, 64)

	if valueTag > math.MaxInt32 {
		return &MaxEnumValueTagError{
			File: unit.File,
			Line: line,
			Col:  col,
		}
	}

	prefix := fullyQualifyIdentifier(scope, "")

	if decl, found := sym.SearchDecl(prefix); found {
		if _, ok := decl.Fields[valueName]; ok {
			return &FieldNameReusedError{
				ParentName: prefix,
				Name:       valueName,
				File:       unit.File,
				Line:       line,
				Col:        col,
			}
		}
		if _, ok := decl.FieldTags[valueTag]; ok {
			return &FieldTagReusedError{
				ParentName: prefix,
				Tag:        valueTag,
				File:       unit.File,
				Line:       line,
				Col:        col,
			}
		}

		ref := symtab.Ref{Unit: unit, Line: line, Col: col, Tag: valueTag}
		sym.AddRef(&decl, prefix, valueName, valueTag, ref)
	} else {
		panic("it should never happen! we should be in the right scope")
	}

	return nil
}
