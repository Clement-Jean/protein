package typecheck

import (
	"fmt"
	"strings"
	"unique"
)

func (tc *TypeChecker) handleEnum(multiset *typeMultiset, pkg []string, unit *Unit, idx uint32) {
	idx += 1

	start := unit.Toks.TokenInfos[idx].Offset
	end := unit.Toks.TokenInfos[idx+1].Offset
	name := strings.TrimSpace(string(unit.Buffer.Range(start, end)))
	prefix := strings.Join(pkg, ".")

	if len(prefix) != 0 && !strings.HasPrefix(prefix, ".") {
		prefix = "." + prefix
	}

	multiset.offsets = append(multiset.offsets, start)
	multiset.names = append(multiset.names, unique.Make(fmt.Sprintf("%s.%s", prefix, name)))
}
