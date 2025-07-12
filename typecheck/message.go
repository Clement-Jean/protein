package typecheck

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Clement-Jean/protein/lexer"
	"github.com/Clement-Jean/protein/parser"
	"github.com/Clement-Jean/protein/symtab"
	"github.com/Clement-Jean/protein/unit"
)

const maxFieldTag = 536870911

func collectIdentifier(idx uint32, unit *unit.Unit, start lexer.TokenInfo) (uint32, string) {
	var name strings.Builder

	for start.Kind == lexer.TokenKindIdentifier || start.Kind == lexer.TokenKindDot {
		end := unit.Toks.TokenInfos[idx+1]
		part := string(unit.Buffer.Range(start.Offset, end.Offset))

		// FIX: this is a hack for avoiding reading the field name
		if strings.HasSuffix(part, " ") {
			name.WriteString(strings.TrimSpace(part))
			break
		}

		name.WriteString(part)
		idx++
		start = unit.Toks.TokenInfos[idx]
	}

	return idx + 1, name.String()
}

func fullyQualifyIdentifier(scope []string, id string) string {
	var fullyQualified string
	if len(scope) != 0 {
		if strings.HasPrefix(scope[0], ".") {
			fullyQualified = strings.Join(scope, ".")
		} else {
			fullyQualified = "." + strings.Join(scope, ".")
		}
	} else if !strings.HasPrefix(id, ".") {
		fullyQualified = "."
	}
	return fullyQualified
}

func (tc *TypeChecker) handleMessage(sym symtab.Symtab, scope *[]string, unit *unit.Unit, idx uint32) error {
	start := unit.Toks.TokenInfos[idx].Offset
	end := unit.Toks.TokenInfos[idx+1].Offset
	line, col := tc.getLineColumn(unit, start)
	name := strings.TrimSpace(string(unit.Buffer.Range(start, end)))
	prefix := strings.Join(*scope, ".")

	if len(prefix) != 0 && !strings.HasPrefix(prefix, ".") {
		prefix = "." + prefix
	}

	fullName := fmt.Sprintf("%s.%s", prefix, name)

	if decl, ok := sym[fullName]; ok {
		return &TypeRedefinedError{
			Name:  fullName,
			Files: []string{unit.File, decl.Unit.File},
			Lines: []int{line, decl.Line},
			Cols:  []int{col, decl.Col},
		}
	}

	(*scope) = append((*scope), name)
	sym[fullName] = symtab.Decl{
		Unit:      unit,
		Line:      line,
		Col:       col,
		Type:      parser.NodeKindMessageDecl,
		Name:      name,
		Fields:    make(map[string]symtab.Ref),
		FieldTags: make(map[int64]symtab.Ref),
	}
	return nil
}

func (tc *TypeChecker) handleOneof(sym symtab.Symtab, scope []string, unit *unit.Unit, idx uint32) error {
	start := unit.Toks.TokenInfos[idx].Offset
	end := unit.Toks.TokenInfos[idx+1].Offset
	line, col := tc.getLineColumn(unit, start)
	name := strings.TrimSpace(string(unit.Buffer.Range(start, end)))
	prefix := strings.Join(scope, ".")

	if len(prefix) != 0 && !strings.HasPrefix(prefix, ".") {
		prefix = "." + prefix
	}

	fullName := fmt.Sprintf("%s.%s", prefix, name)

	if decl, ok := sym[fullName]; ok {
		return &TypeRedefinedError{
			Name:  fullName,
			Files: []string{unit.File, decl.Unit.File},
			Lines: []int{line, decl.Line},
			Cols:  []int{col, decl.Col},
		}
	}

	sym[fullName] = symtab.Decl{
		Unit: unit,
		Line: line,
		Col:  col,
		Type: parser.NodeKindOneOfDecl,
		Name: name,
	}
	return nil
}

func (tc *TypeChecker) handleMapDecl(sym symtab.Symtab, scope []string, unit *unit.Unit, idx uint32) error {
	start := unit.Toks.TokenInfos[idx]
	endIdx, fieldName := collectIdentifier(idx, unit, start)
	line, col := tc.getLineColumn(unit, start.Offset)

	endOffset := unit.Toks.TokenInfos[endIdx+1].Offset
	endNextOffset := unit.Toks.TokenInfos[endIdx+2].Offset
	fieldTag, _ := strconv.ParseInt(string(unit.Buffer.Range(endOffset, endNextOffset)), 10, 32)

	if fieldTag > maxFieldTag {
		return &MaxFieldTagError{
			File: unit.File,
			Line: line,
			Col:  col,
		}
	}

	prefix := fullyQualifyIdentifier(scope, "")

	if value, found := sym[prefix]; found {
		if _, ok := value.Fields[fieldName]; ok {
			return &FieldNameReusedError{
				ParentName: prefix,
				Name:       fieldName,
				File:       unit.File,
				Line:       line,
				Col:        col,
			}
		}
		if _, ok := value.FieldTags[fieldTag]; ok {
			return &FieldTagReusedError{
				ParentName: prefix,
				Tag:        fieldTag,
				File:       unit.File,
				Line:       line,
				Col:        col,
			}
		}

		ref := symtab.Ref{Unit: unit, Line: line, Col: col, Tag: fieldTag}
		value.Fields[fieldName] = ref
		value.FieldTags[fieldTag] = ref
	} else {
		panic("it should never happen! we should be in the right scope")
	}

	return nil
}

func (tc *TypeChecker) handleMapValue(sym symtab.Symtab, scope []string, unit *unit.Unit, idx uint32) error {
	start := unit.Toks.TokenInfos[idx]
	endIdx, id := collectIdentifier(idx, unit, start)

	if len(id) == 0 { // non user-defined types, no need to check
		return nil
	}

	line, col := tc.getLineColumn(unit, start.Offset)
	isPrecededByDot := idx-1 > 0 && unit.Toks.TokenInfos[idx-1].Kind == lexer.TokenKindDot
	endOffset := unit.Toks.TokenInfos[idx].Offset
	endNextOffset := unit.Toks.TokenInfos[endIdx-1].Offset
	fieldName := strings.TrimSuffix(string(unit.Buffer.Range(endOffset, endNextOffset)), " ")

	prefix := fullyQualifyIdentifier(scope, id)
	name := splitAndMerge(id, prefix)

	if isPrecededByDot {
		start = unit.Toks.TokenInfos[idx-1]
	}

	if value, found := sym[prefix]; found {
		ref := symtab.Ref{
			Unit:     unit,
			Line:     line,
			Col:      col,
			TypeName: name,
		}
		value.Fields[fieldName] = ref
	} else {
		panic("it should never happen! we should be in the right scope")
	}
	return nil
}

func (tc *TypeChecker) handleField(sym symtab.Symtab, scope []string, unit *unit.Unit, idx uint32) error {
	start := unit.Toks.TokenInfos[idx]
	isPrecededByDot := idx-1 > 0 && unit.Toks.TokenInfos[idx-1].Kind == lexer.TokenKindDot
	endIdx, id := collectIdentifier(idx, unit, start)
	line, col := tc.getLineColumn(unit, start.Offset)

	endOffset := unit.Toks.TokenInfos[endIdx].Offset
	endNextOffset := unit.Toks.TokenInfos[endIdx+1].Offset
	fieldName := strings.TrimSuffix(string(unit.Buffer.Range(endOffset, endNextOffset)), " ")

	endOffset = unit.Toks.TokenInfos[endIdx+2].Offset
	endNextOffset = unit.Toks.TokenInfos[endIdx+3].Offset
	fieldTag, _ := strconv.ParseInt(string(unit.Buffer.Range(endOffset, endNextOffset)), 10, 32)

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

	if value, found := sym[prefix]; found {
		if _, ok := value.Fields[fieldName]; ok {
			return &FieldNameReusedError{
				ParentName: prefix,
				Name:       fieldName,
				File:       unit.File,
				Line:       line,
				Col:        col,
			}
		}
		if _, ok := value.FieldTags[fieldTag]; ok {
			return &FieldTagReusedError{
				ParentName: prefix,
				Tag:        fieldTag,
				File:       unit.File,
				Line:       line,
				Col:        col,
			}
		}

		var ref symtab.Ref

		if len(id) == 0 { // non user-defined types (e.g. int32)
			ref = symtab.Ref{Unit: unit, Line: line, Col: col, Type: value.Type, Tag: fieldTag}
		} else {
			ref = symtab.Ref{Unit: unit, Line: line, Col: col, TypeName: name, Tag: fieldTag}
		}

		value.Fields[fieldName] = ref
		value.FieldTags[fieldTag] = ref
	} else {
		panic("it should never happen! we should be in the right scope")
	}

	return nil
}
