package linker

func partition(a []string, lo, hi, d int) (lt, gt, p int) {
	var pivot int

	if d >= len(a[lo]) {
		pivot = -1
	} else {
		pivot = int(a[lo][d])
	}
	i := lo + 1

	for i <= hi {
		var t int
		if d >= len(a[i]) {
			t = -1
		} else {
			t = int(a[i][d])
		}

		if t < pivot {
			a[lo], a[i] = a[i], a[lo]
			lo++
			i++
		} else if t > pivot {
			a[i], a[hi] = a[hi], a[i]
			hi--
		} else {
			i++
		}
	}

	return lo, hi, pivot
}

func multisetSort(a []string) []bool {
	unique := make([]bool, len(a))
	var q [][3]int

	q = append(q, [3]int{0, len(a) - 1, 0})

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

		unique[lt] = true

		q = append(q, [3]int{lo, lt - 1, d})
		if pivot >= 0 {
			q = append(q, [3]int{lt, gt, d + 1})
		}
		q = append(q, [3]int{gt + 1, hi, d})
	}

	return unique
}
