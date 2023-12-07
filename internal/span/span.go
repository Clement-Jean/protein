package span

import "fmt"

type Span struct {
	Start int
	End   int
}

func (s Span) String() string {
	return fmt.Sprintf("{ start: %d, end: %d }", s.Start, s.End)
}
