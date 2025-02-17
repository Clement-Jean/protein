package linker

import (
	"fmt"
	"strings"
	"unique"
)

func (l *Linker) handleService(multiset *typeMultiset, pkg []string, unit Unit, idx uint32) {
	idx += 1

	start := unit.Toks.TokenInfos[idx].Offset
	end := unit.Toks.TokenInfos[idx+1].Offset
	name := strings.TrimSpace(string(unit.Buffer.Range(start, end)))

	prefix := strings.Join(pkg, ".")
	(*multiset).names = append((*multiset).names, unique.Make(fmt.Sprintf("%s.%s", prefix, name)))
}
