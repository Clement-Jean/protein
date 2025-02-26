package typecheck

import (
	"fmt"
	"strings"
	"unique"

	"github.com/Clement-Jean/protein/lexer"
)

func collectIdentifier(idx uint32, unit *Unit, start lexer.TokenInfo) string {
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

	return name.String()
}

func (tc *TypeChecker) handleMessage(multiset *typeMultiset, pkg *[]string, unit *Unit, idx uint32) {
	idx += 1

	start := unit.Toks.TokenInfos[idx].Offset
	end := unit.Toks.TokenInfos[idx+1].Offset
	name := strings.TrimSpace(string(unit.Buffer.Range(start, end)))
	prefix := strings.Join(*pkg, ".")

	multiset.offsets = append(multiset.offsets, start)
	multiset.names = append(multiset.names, unique.Make(fmt.Sprintf("%s.%s", prefix, name)))
	(*pkg) = append((*pkg), name)
}

func (tc *TypeChecker) handleOneof(multiset *typeMultiset, pkg []string, unit *Unit, idx uint32) {
	idx += 1

	start := unit.Toks.TokenInfos[idx].Offset
	end := unit.Toks.TokenInfos[idx+1].Offset
	name := strings.TrimSpace(string(unit.Buffer.Range(start, end)))
	prefix := strings.Join(pkg, ".")

	multiset.offsets = append(multiset.offsets, start)
	multiset.names = append(multiset.names, unique.Make(fmt.Sprintf("%s.%s", prefix, name)))
}

func (tc *TypeChecker) handleMapValue(multiset *typeMultiset, pkg []string, unit *Unit, idx uint32) {
	start := unit.Toks.TokenInfos[idx]
	isPrecededByDot := idx-1 > 0 && unit.Toks.TokenInfos[idx-1].Kind == lexer.TokenKindDot
	id := collectIdentifier(idx, unit, start)

	if len(id) == 0 { // non user-defined types (e.g. int32)
		return
	}

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

func (tc *TypeChecker) handleField(multiset *typeMultiset, pkg []string, unit *Unit, idx uint32) {
	start := unit.Toks.TokenInfos[idx]
	isPrecededByDot := idx-1 > 0 && unit.Toks.TokenInfos[idx-1].Kind == lexer.TokenKindDot
	id := collectIdentifier(idx, unit, start)

	if len(id) == 0 { // non user-defined types (e.g. int32)
		return
	}

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
