package typecheck

import (
	"fmt"
	"strings"

	"github.com/Clement-Jean/protein/lexer"
	"github.com/Clement-Jean/protein/parser"
	"github.com/Clement-Jean/protein/symtab"
	"github.com/Clement-Jean/protein/unit"
)

func (tc *TypeChecker) handleExtend(sym *symtab.Symtab, scope *[]string, unit *unit.Unit, idx uint32) error {
	start := unit.Toks.TokenInfos[idx]
	_, name := collectIdentifier(idx, unit, start)
	line, col := tc.getLineColumn(unit, start.Offset)

	if !strings.HasPrefix(name, ".") {
		name = "." + name
	}

	decl, ok := sym.SearchDecl(name)
	if !ok { // message decl not defined yet
		decl = symtab.Decl{
			Unit:      unit,
			Line:      line,
			Col:       col,
			Type:      parser.NodeKindExtendDecl,
			Name:      name,
			Fields:    make(map[string]symtab.Ref),
			FieldTags: make(map[int64]symtab.Ref),
		}
	}

	// FIX is the ref really needed?
	ref := symtab.Ref{
		Unit:     unit,
		Line:     line,
		Col:      col,
		TypeName: name,
	}
	sym.AddRef(&decl, "", 0, ref)
	sym.Decls[name] = decl
	return nil
}

func (tc *TypeChecker) handleExtendField(sym *symtab.Symtab, scope []string, unit *unit.Unit, idx uint32) error {
	start := unit.Toks.TokenInfos[idx]
	isPrecededByDot := idx-1 > 0 && unit.Toks.TokenInfos[idx-1].Kind == lexer.TokenKindDot
	endIdx, id := collectIdentifier(idx, unit, start)
	line, col := tc.getLineColumn(unit, start.Offset)

	endOffset := unit.Toks.TokenInfos[endIdx].Offset
	endNextOffset := unit.Toks.TokenInfos[endIdx+1].Offset
	fieldName := strings.TrimSuffix(string(unit.Buffer.Range(endOffset, endNextOffset)), " ")

	tagToken := unit.Toks.TokenInfos[endIdx+2]
	endOffset = tagToken.Offset
	endNextOffset = unit.Toks.TokenInfos[endIdx+3].Offset
	tag := strings.TrimRight(string(unit.Buffer.Range(endOffset, endNextOffset)), "\n ")
	fieldTag := parseTag(tagToken, tag)

	if fieldTag > maxFieldTag {
		return &MaxFieldTagError{
			File: unit.File,
			Line: line,
			Col:  col,
		}
	}

	prefix := fullyQualifyIdentifier(scope, id)
	name := splitAndMerge(id, prefix)

	if isPrecededByDot {
		start = unit.Toks.TokenInfos[idx-1]
	}

	declName := fmt.Sprintf("%s.%s", prefix, fieldName)

	if found, ok := sym.SearchDecl(declName); ok {
		return &TypeRedefinedError{
			Files: []string{unit.File, found.Unit.File},
			Lines: []int{line, found.Line},
			Cols:  []int{col, found.Col},
			Name:  declName,
		}
	}

	decl := symtab.Decl{
		Unit:      unit,
		Line:      line,
		Col:       col,
		Type:      parser.NodeKindExtendFieldDecl,
		Name:      declName,
		Fields:    make(map[string]symtab.Ref),
		FieldTags: make(map[int64]symtab.Ref),
	}

	var ref symtab.Ref
	if len(id) == 0 { // non user-defined types (e.g. int32)
		ref = symtab.Ref{Unit: unit, Line: line, Col: col, Type: parser.NodeKindUndefined, Tag: fieldTag}
	} else {
		ref = symtab.Ref{Unit: unit, Line: line, Col: col, TypeName: name, Tag: fieldTag}
	}

	sym.AddRef(&decl, fieldName, fieldTag, ref)
	sym.Decls[declName] = decl
	return nil
}
