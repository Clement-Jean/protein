package typecheck

import (
	"strings"

	"github.com/Clement-Jean/protein/lexer"
)

func (tc *TypeChecker) handlePackage(pkgs map[*Unit]string, unit *Unit, idx uint32) {
	idx += 1

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
