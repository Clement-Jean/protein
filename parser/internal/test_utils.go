package internal

import (
	"strings"

	"github.com/Clement-Jean/protein/internal/span"
)

func CaseName(domain string, isValid bool, descriptions ...string) string {
	sb := strings.Builder{}

	sb.WriteString(domain)

	if isValid {
		sb.WriteString("_valid")
	} else {
		sb.WriteString("_invalid")
	}

	for _, description := range descriptions {
		sb.WriteString("_")
		sb.WriteString(description)
	}

	return sb.String()
}

// ReferenceString creates a map out of indices.
// it checks for characters that are different from sep
// and record their location in a map.
//
// e.g a--------b-----c gives:
// a: 0
// b: 9
// c: 15
func ReferenceString(indices string, sep rune) map[rune]int {
	ref := map[rune]int{}
	column := 0

	for _, index := range indices {
		if index != sep && index != '\n' {
			ref[index] = column
		}

		column++
	}

	return ref
}

func MakeSpansFromIndices(ref map[rune]int, indices [][2]rune) (spans []span.Span) {
	for _, index := range indices {
		start := ref[index[0]]
		end := ref[index[1]]

		if end < start {
			panic("end shouldn't be bigger than start")
		}

		spans = append(spans, span.Span{Start: ref[index[0]], End: ref[index[1]]})
	}

	// Add EOF span
	last := spans[len(spans)-1]
	spans = append(spans, span.Span{Start: last.End, End: last.End})
	return spans
}

func EmptyErrorSliceIfNil(item error) []error {
	if item == nil {
		return nil
	}

	return []error{item}
}
