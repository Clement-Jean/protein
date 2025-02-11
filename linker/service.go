package linker

import (
	"fmt"
	"strings"
)

func (l *Linker) handleService(multiset *[]string, pkg []string, unit Unit, idx uint32) {
	idx += 1

	start := unit.Toks.TokenInfos[idx].Offset
	end := unit.Toks.TokenInfos[idx+1].Offset
	name := strings.TrimSpace(string(unit.Buffer.Range(start, end)))

	prefix := strings.Join(pkg, ".")
	(*multiset) = append((*multiset), fmt.Sprintf("%s.%s", prefix, name))
}
