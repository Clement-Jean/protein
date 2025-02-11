package linker

import (
	"strings"

	"github.com/Clement-Jean/protein/lexer"
)

func (l *Linker) handlePackage(pkgs *map[string]string, unit Unit, idx uint32) {
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

	(*pkgs)[unit.File] = name.String()
}
