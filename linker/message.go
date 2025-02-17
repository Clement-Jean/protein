package linker

import (
	"fmt"
	"strings"
	"unique"

	"github.com/Clement-Jean/protein/lexer"
)

func collectIdentifier(idx uint32, unit Unit, start lexer.TokenInfo) string {
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

func (l *Linker) handleMessage(multiset *typeMultiset, pkg *[]string, unit Unit, idx uint32) {
	idx += 1

	start := unit.Toks.TokenInfos[idx].Offset
	end := unit.Toks.TokenInfos[idx+1].Offset
	name := strings.TrimSpace(string(unit.Buffer.Range(start, end)))
	prefix := strings.Join(*pkg, ".")

	(*multiset).names = append((*multiset).names, unique.Make(fmt.Sprintf("%s.%s", prefix, name)))
	(*pkg) = append((*pkg), name)
}

func (l *Linker) handleMapValue(multiset *typeMultiset, pkg []string, unit Unit, idx uint32) {
	start := unit.Toks.TokenInfos[idx]

	if start.Kind == lexer.TokenKindDot {
		start = unit.Toks.TokenInfos[idx+1]
	}

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

	isFullyQualified := (idx-1 > 0 && unit.Toks.TokenInfos[idx-1].Kind == lexer.TokenKindDot) || strings.Contains(id, ".")

	if !isFullyQualified {
		prefix := strings.Join(pkg, ".")
		(*multiset).names = append((*multiset).names, unique.Make(fmt.Sprintf("%s.%s", prefix, id)))
		return
	}

	(*multiset).names = append((*multiset).names, unique.Make(id))
}

func (l *Linker) handleField(multiset *typeMultiset, pkg []string, unit Unit, idx uint32) {
	switch unit.Toks.TokenInfos[idx].Kind {
	case lexer.TokenKindOptional, lexer.TokenKindRequired, lexer.TokenKindRepeated:
		idx += 1
	}

	start := unit.Toks.TokenInfos[idx]
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

	isFullyQualified := (idx-1 > 0 && unit.Toks.TokenInfos[idx-1].Kind == lexer.TokenKindDot) || strings.Contains(id, ".")

	if !isFullyQualified {
		prefix := strings.Join(pkg, ".")
		(*multiset).names = append((*multiset).names, unique.Make(fmt.Sprintf("%s.%s", prefix, id)))
		return
	}

	(*multiset).names = append((*multiset).names, unique.Make(id))
}
