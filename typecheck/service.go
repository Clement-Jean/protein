package typecheck

import (
	"strings"
	"unique"

	"github.com/Clement-Jean/protein/lexer"
)

func (tc *TypeChecker) handleService(multiset *typeMultiset, pkg []string, unit *Unit, idx uint32) {
	// idx += 1

	// start := unit.Toks.TokenInfos[idx].Offset
	// end := unit.Toks.TokenInfos[idx+1].Offset
	// name := strings.TrimSpace(string(unit.Buffer.Range(start, end)))
	// prefix := strings.Join(pkg, ".")

	// multiset.offsets = append(multiset.offsets, start)
	// multiset.names = append(multiset.names, unique.Make(fmt.Sprintf("%s.%s", prefix, name)))
}

func (tc *TypeChecker) handleRPC(multiset *typeMultiset, pkg []string, unit *Unit, idx uint32) {
	start := unit.Toks.TokenInfos[idx]
	isPrecededByDot := idx-1 > 0 && unit.Toks.TokenInfos[idx-1].Kind == lexer.TokenKindDot
	id := collectIdentifier(idx, unit, start)

	var prefix string
	if len(pkg) != 0 {
		prefix = strings.Join(pkg, ".")
	} else {
		prefix = "."
	}

	name := splitAndMerge(id, prefix)

	if isPrecededByDot {
		start = unit.Toks.TokenInfos[idx-1]
	}

	multiset.offsets = append(multiset.offsets, start.Offset)
	multiset.names = append(multiset.names, unique.Make(name))
}
