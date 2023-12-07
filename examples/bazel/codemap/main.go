package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/Clement-Jean/protein/codemap"
	"github.com/Clement-Jean/protein/lexer"
)

func main() {
	content := "message Example {}"

	cm := codemap.New()
	fm := cm.Insert("example.proto", []byte(content))
	l := lexer.New([]byte(content))
	kinds, spans := l.Tokenize()
	tokens := fm.RegisterTokens(kinds, spans)

	w := tabwriter.NewWriter(os.Stdout, 0, 8, 1, ' ', 0)

	fmt.Fprintln(w, "Token\tLiteral")
	for _, token := range tokens {
		literal := fm.Lookup(token.ID)
		fmt.Fprintf(w, "%v\t%s\n", token, literal)
	}

	w.Flush()
}
