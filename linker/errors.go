package linker

import (
	"fmt"
	"strings"
)

type ImportCycleError struct {
	Files []string
}

func (e *ImportCycleError) Error() string {
	var msg strings.Builder

	for _, file := range e.Files {
		msg.WriteString(file)
		msg.WriteString(" -> ")
	}
	msg.WriteString(e.Files[0])

	return fmt.Sprintf("cycle found: %s", msg.String())
}
