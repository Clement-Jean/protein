package typecheck

import (
	"strings"

	"github.com/Clement-Jean/protein/lexer"
	"github.com/Clement-Jean/protein/parser"
	"github.com/Clement-Jean/protein/symtab"
	"github.com/Clement-Jean/protein/unit"
)

func (tc *TypeChecker) registerPackage(pkgs map[*unit.Unit]string, unit *unit.Unit, idx uint32) {
	var name strings.Builder

	start := unit.Toks.TokenInfos[idx]
	for start.Kind == lexer.TokenKindIdentifier || start.Kind == lexer.TokenKindDot {
		end := unit.Toks.TokenInfos[idx+1]
		part := string(unit.Buffer.Range(start.Offset, end.Offset))

		name.WriteString(part)
		idx++
		start = unit.Toks.TokenInfos[idx]
	}

	pkgs[unit] = name.String()
}

func (tc *TypeChecker) handlePackage(sym *symtab.Symtab, unit *unit.Unit, idx uint32) error {
	start := unit.Toks.TokenInfos[idx]
	_, pkgName := collectIdentifier(idx, unit, start)
	line, col := tc.getLineColumn(unit, start.Offset)
	parts := strings.Split(pkgName, ".")

	for i := len(parts) - 1; i >= 0; i-- {
		name := "." + strings.Join(parts[:i+1], ".")

		if decl, ok := sym.Decls[name]; ok && decl.Type != parser.NodeKindPackageStmt { // message is defined already
			return &TypeRedefinedError{
				Files: []string{decl.Unit.File, unit.File},
				Lines: []int{decl.Line, line},
				Cols:  []int{decl.Col, col},
				Name:  name,
			}
		}

		sym.Decls[name] = symtab.Decl{
			Unit: unit,
			Line: line,
			Col:  col,
			Type: parser.NodeKindPackageStmt,
		}
	}
	return nil
}
