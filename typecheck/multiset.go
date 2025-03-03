package typecheck

import (
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
		pivot = int(a.names[lo].Value()[d])
	}
	i := lo + 1

	for i <= hi {
		if a.names[i] == a.names[lo] {
			if a.kinds[i] < a.kinds[lo] {
				a.kinds[i], a.kinds[lo] = a.kinds[lo], a.kinds[i]
				a.offsets[i], a.offsets[lo] = a.offsets[lo], a.offsets[i]
				a.units[i], a.units[lo] = a.units[lo], a.units[i]
			} else if a.kinds[i] > a.kinds[lo] {
				a.kinds[i], a.kinds[hi] = a.kinds[hi], a.kinds[i]
				a.offsets[i], a.offsets[hi] = a.offsets[hi], a.offsets[i]
				a.units[i], a.units[hi] = a.units[hi], a.units[i]
			}

			i++
			continue
		}

		var t int
		if d >= len(a.names[i].Value()) {
			t = -1
		} else {
			t = int(a.names[i].Value()[d])
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
