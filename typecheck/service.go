package typecheck

import (
	"fmt"
	"strings"

	"github.com/Clement-Jean/protein/lexer"
	"github.com/Clement-Jean/protein/parser"
	"github.com/Clement-Jean/protein/symtab"
	"github.com/Clement-Jean/protein/unit"
)

func (tc *TypeChecker) handleService(sym *symtab.Symtab, scope *[]string, unit *unit.Unit, idx uint32) error {
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
		Unit:   unit,
		Line:   line,
		Col:    col,
		Type:   parser.NodeKindServiceDecl,
		Name:   name,
		Fields: make(map[string]symtab.Ref),
	}
	return nil
}

func (tc *TypeChecker) handleRPCDecl(sym *symtab.Symtab, scope []string, unit *unit.Unit, idx uint32) error {
	start := unit.Toks.TokenInfos[idx]
	_, rpcName := collectIdentifier(idx, unit, start)
	line, col := tc.getLineColumn(unit, start.Offset)
	prefix := strings.Join(scope, ".")

	if len(prefix) != 0 && !strings.HasPrefix(prefix, ".") {
		prefix = "." + prefix
	}

	fullName := fmt.Sprintf("%s.%s", prefix, rpcName)

	if decl, ok := sym.SearchDecl(fullName); ok {
		return &TypeRedefinedError{
			Name:  fullName,
			Files: []string{unit.File, decl.Unit.File},
			Lines: []int{line, decl.Line},
			Cols:  []int{col, decl.Col},
		}
	}

	sym.Decls[fullName] = symtab.Decl{
		Unit:   unit,
		Line:   line,
		Col:    col,
		Type:   parser.NodeKindRPCDecl,
		Name:   rpcName,
		Fields: make(map[string]symtab.Ref),
	}
	return nil
}

func (tc *TypeChecker) handleRPCInputOutput(sym *symtab.Symtab, scope []string, unit *unit.Unit, idx uint32) error {
	start := unit.Toks.TokenInfos[idx]
	isPrecededByDot := idx-1 > 0 && unit.Toks.TokenInfos[idx-1].Kind == lexer.TokenKindDot
	_, id := collectIdentifier(idx, unit, start)
	line, col := tc.getLineColumn(unit, start.Offset)

	prefix := fullyQualifyIdentifier(scope, id)
	name := splitAndMerge(id, prefix)

	if isPrecededByDot {
		start = unit.Toks.TokenInfos[idx-1]
	}

	if decl, found := sym.SearchDecl(prefix); found {
		ref := symtab.Ref{
			Unit:     unit,
			Line:     line,
			Col:      col,
			Type:     parser.NodeKindRPCInputOutput,
			TypeName: name,
		}
		sym.AddRef(&decl, name, 0, ref)
	} else {
		panic("it should never happen! we should be in the right scope")
	}
	return nil
}
