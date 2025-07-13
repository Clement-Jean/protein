package typecheck

import (
	"fmt"
	"strings"

	"github.com/Clement-Jean/protein/symtab"
)

func checkUpperScopes(sym *symtab.Symtab, typeName string) (string, symtab.Decl, bool) {
	idxEnd := strings.IndexByte(typeName, ']')
	if idxEnd == -1 {
		return typeName, symtab.Decl{}, false
	}

	idxStart := strings.IndexByte(typeName, '[')
	if idxStart == -1 {
		idxStart = 0
	}

	minScope := typeName[0:idxStart]
	scope := typeName[idxStart+1 : idxEnd]
	ref := typeName[idxEnd+1:]

	pkgName := scope
	name := fmt.Sprintf("%s.%s.%s", minScope, pkgName, ref)
	scopeIdx := strings.LastIndexByte(pkgName, '.')

	for {
		if decl, ok := sym.SearchDecl(name); ok {
			return name, decl, ok
		}

		scopeIdx = strings.LastIndexByte(pkgName, '.')
		if scopeIdx == -1 {
			break
		}

		pkgName = pkgName[:scopeIdx]
		name = fmt.Sprintf("%s.%s.%s", minScope, pkgName, ref)
	}

	name = fmt.Sprintf("%s.%s", minScope, ref)
	decl, ok := sym.SearchDecl(name)
	return name, decl, ok
}

func splitAndMerge(id, scope string) string {
	if len(id) < 1 || id[0] == '.' { // fully qualified
		return id
	} else if len(scope) < 1 || scope == "." { // global scope
		return "." + id
	}

	r := id
	s := scope
	if s[0] == '.' {
		s = s[1:]
	}

	var sb strings.Builder

	sb.Grow(len(id) + len(scope) + 2)

	lastEqual := -1
	idxRef := strings.IndexByte(r, '.')
	idxScope := strings.IndexByte(s, '.')

	var lastScope string
	var lastRef string
	if idxRef != -1 {
		lastRef = r[:idxRef]
		r = r[idxRef+1:]
	}

	for idxScope != -1 && idxRef != -1 {
		lastScope = s[:idxScope]

		if idxScope >= len(s) {
			return id // FIX: error handling
		}

		s = s[idxScope+1:]

		if len(lastScope) == 0 {
			continue
		}

		sb.WriteByte('.')
		if lastScope == lastRef {
			if lastEqual == -1 {
				lastEqual = sb.Len() - 1
			}

			idxRef = strings.IndexByte(r, '.')
			if idxRef != -1 {
				lastRef = r[:idxRef]
				r = r[idxRef+1:]
			}
		}
		sb.WriteString(lastScope)

		idxScope = strings.IndexByte(s, '.')
	}

	if s == lastRef {
		sb.WriteByte('.')
		if lastEqual == -1 {
			lastEqual = sb.Len() - 1
		}
		sb.WriteString(lastRef)
	} else if lastScope == r {
		sb.WriteByte('.')
		if lastEqual == -1 {
			lastEqual = sb.Len() - 1
		}
		sb.WriteString(lastScope)
	} else if len(s) != 0 {
		sb.WriteByte('.')
		sb.WriteString(s)
	}

	if len(r) != 0 {
		if len(s) != 0 {
			sb.WriteByte(']')
		}
		sb.WriteString(id)
	}

	if lastEqual <= 0 {
		return sb.String()
	}

	b := []byte(sb.String())
	b[lastEqual] = '['
	return string(b)
}
