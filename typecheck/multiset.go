package typecheck

import (
	"cmp"
	"fmt"
	"slices"
	"strings"
	"unique"

	"github.com/Clement-Jean/protein/parser"
)

type typeMultiset struct {
	names   []unique.Handle[string]
	kinds   []parser.NodeKind
	offsets []uint32
	units   []*Unit
}

func partition(a *typeMultiset, lo, hi, d int) (lt, gt, pivot int) {
	if d >= len(a.names[lo].Value()) {
		pivot = -1
	} else {
		switch a.names[lo].Value()[d] {
		case '[', ']':
			pivot = int('.')
		default:
			pivot = int(a.names[lo].Value()[d])
		}
	}
	i := lo + 1

	for i <= hi {
		if a.names[i] == a.names[lo] {
			if a.kinds[i] < a.kinds[lo] {
				a.kinds[lo], a.kinds[i] = a.kinds[i], a.kinds[lo]
				a.offsets[lo], a.offsets[i] = a.offsets[i], a.offsets[lo]
				a.units[lo], a.units[i] = a.units[i], a.units[lo]
				lo++
				i++
			} else if a.kinds[i] > a.kinds[lo] {
				a.kinds[i], a.kinds[hi] = a.kinds[hi], a.kinds[i]
				a.offsets[i], a.offsets[hi] = a.offsets[hi], a.offsets[i]
				a.units[i], a.units[hi] = a.units[hi], a.units[i]
				hi--
			} else {
				i++
			}
			continue
		}

		var t int
		if d >= len(a.names[i].Value()) {
			t = -1
		} else {
			switch a.names[i].Value()[d] {
			case '[', ']':
				t = int('.')
			default:
				t = int(a.names[i].Value()[d])
			}
		}

		if t < pivot {
			a.names[lo], a.names[i] = a.names[i], a.names[lo]
			a.kinds[lo], a.kinds[i] = a.kinds[i], a.kinds[lo]
			a.offsets[lo], a.offsets[i] = a.offsets[i], a.offsets[lo]
			a.units[lo], a.units[i] = a.units[i], a.units[lo]
			lo++
			i++
		} else if t > pivot {
			a.names[i], a.names[hi] = a.names[hi], a.names[i]
			a.kinds[i], a.kinds[hi] = a.kinds[hi], a.kinds[i]
			a.offsets[i], a.offsets[hi] = a.offsets[hi], a.offsets[i]
			a.units[i], a.units[hi] = a.units[hi], a.units[i]
			hi--
		} else {
			i++
		}
	}

	return lo, hi, pivot
}

func multisetSort(a *typeMultiset) {
	n := len(a.names)
	var q [][3]int

	q = append(q, [3]int{0, n - 1, 0})

	for len(q) != 0 {
		front := q[0]
		lo := front[0]
		hi := front[1]
		d := front[2]

		q = q[1:]

		if hi <= lo {
			continue
		}

		lt, gt, pivot := partition(a, lo, hi, d)

		q = append(q, [3]int{lo, lt - 1, d})
		if pivot >= 0 {
			q = append(q, [3]int{lt, gt, d + 1})
		}
		q = append(q, [3]int{gt + 1, hi, d})
	}
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

func checkUpperScopes(decls *typeMultiset, s string) (int, string, bool) {
	cmpFn := func(h unique.Handle[string], s string) int {
		return cmp.Compare(h.Value(), s)
	}

	idxEnd := strings.IndexByte(s, ']')
	if idxEnd == -1 {
		//println("check", s)
		if idx, ok := slices.BinarySearchFunc(decls.names, s, cmpFn); ok {
			return idx, s, ok
		}
		return -1, s, false
	}

	idxStart := strings.IndexByte(s, '[')
	if idxStart == -1 {
		idxStart = 0
	}

	minScope := s[0:idxStart]
	scope := s[idxStart+1 : idxEnd]
	ref := s[idxEnd+1:]

	pkgName := scope
	name := fmt.Sprintf("%s.%s.%s", minScope, pkgName, ref)
	scopeIdx := strings.LastIndexByte(pkgName, '.')

	for {
		//println("check", name)
		if idx, ok := slices.BinarySearchFunc(decls.names, name, cmpFn); ok {
			return idx, name, ok
		}

		scopeIdx = strings.LastIndexByte(pkgName, '.')
		if scopeIdx == -1 {
			break
		}

		pkgName = pkgName[:scopeIdx]
		name = fmt.Sprintf("%s.%s.%s", minScope, pkgName, ref)
	}

	name = fmt.Sprintf("%s.%s", minScope, ref)
	//println("check", name)
	idx, ok := slices.BinarySearchFunc(decls.names, name, cmpFn)
	return idx, name, ok
}
