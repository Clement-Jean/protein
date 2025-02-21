package typecheck

import (
	"fmt"
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

	// if name is fully qualified -> add as is
	// else -> add the curr pkg + potential parent messages
	//
	// e.g.
	// google.protobuf.Empty -> google.protobuf.Empty
	// .google.protobuf.Empty -> google.protobuf.Empty
	// D -> the.curr.pkg.Parent.D
	// .D -> D

	isFullyQualified := isPrecededByDot || strings.Contains(id, ".")

	if !isFullyQualified {
		prefix := strings.Join(pkg, ".")
		multiset.offsets = append(multiset.offsets, start.Offset)
		multiset.names = append(multiset.names, unique.Make(fmt.Sprintf("%s.%s", prefix, id)))
		return
	}

	if isPrecededByDot {
		start = unit.Toks.TokenInfos[idx-1]
	}

	multiset.offsets = append(multiset.offsets, start.Offset)
	multiset.names = append(multiset.names, unique.Make(id))
}
