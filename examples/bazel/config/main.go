package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/Clement-Jean/protein/lexer"
)

func main() {
	content := "message Example { /*comment*/ }"

	l := lexer.New([]byte(content))
	kinds, spans := l.Tokenize()
	w := tabwriter.NewWriter(os.Stdout, 0, 8, 1, ' ', 0)

	fmt.Fprintln(w, "Kind\tSpan")
	for i := 0; i < len(kinds); i++ {
		fmt.Fprintf(w, "%v\t%v\n", kinds[i], spans[i])
	}
	w.Flush()
}
