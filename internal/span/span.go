package span

import "fmt"

type Span struct {
	Start uint64
	End   uint64
}

func (s Span) String() string {
	return fmt.Sprintf("{ start: %d, end: %d }", s.Start, s.End)
}
