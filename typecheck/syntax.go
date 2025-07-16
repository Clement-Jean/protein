package typecheck

import (
	"strings"

	"github.com/Clement-Jean/protein/unit"
)

type syntaxKind uint8

const (
	syntaxProto2 syntaxKind = iota
	syntaxProto3
)

func (tc *TypeChecker) handleSyntax(unit *unit.Unit, idx uint32) error {
	start := unit.Toks.TokenInfos[idx].Offset
	end := unit.Toks.TokenInfos[idx+1].Offset
	syntax := strings.Trim(string(unit.Buffer.Range(start, end)), "\"'")

	switch syntax {
	case "proto2":
		tc.unitSyntax[unit] = syntaxProto2
	case "proto3":
		tc.unitSyntax[unit] = syntaxProto3
	default:
		line, col := tc.getLineColumn(unit, start)
		return &UnknownSyntaxError{
			File:  unit.File,
			Line:  line,
			Col:   col,
			Value: syntax,
		}
	}
	return nil
}
