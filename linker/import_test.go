package linker_test

import (
	"strings"
	"testing"

	"github.com/Clement-Jean/protein/lexer"
	"github.com/Clement-Jean/protein/linker"
	"github.com/Clement-Jean/protein/parser"
	"github.com/Clement-Jean/protein/source"
)

func createUnits(t *testing.T, files []string, contents []string) []linker.Unit {
	t.Helper()

	if len(files) != len(contents) {
		t.Fatal("expected the same number of files and contents")
	}

	var units []linker.Unit
	for i, content := range contents {
		s, err := source.NewFromReader(strings.NewReader(content))
		if err != nil {
			panic(err)
		}

		l, err := lexer.NewFromSource(s)
		if err != nil {
			t.Fatal(err)
		}

		tb, errs := l.Lex()
		if len(errs) != 0 {
			t.Fatal(errs)
		}

		p := parser.New(tb)
		pt, errs := p.Parse()
		if len(errs) != 0 {
			t.Fatal(errs)
		}

		units = append(units, linker.Unit{
			File:   files[i],
			Buffer: s,
			Toks:   tb,
			Tree:   pt,
		})
	}

	return units
}

func TestImportCycle(t *testing.T) {
	tests := []struct {
		files    []string
		contents []string
		errors   []string
	}{
		{
			files: []string{"a.proto", "b.proto", "c.proto"},
			contents: []string{
				"import 'b.proto';", // a.proto
				"import 'c.proto';", // b.proto
				"",                  // c.proto
			},
		},
		{
			files: []string{"a.proto", "b.proto"},
			contents: []string{
				"import 'b.proto';", // a.proto
				"import 'a.proto';", // b.proto
			},
			errors: []string{"cycle found: a.proto -> b.proto -> a.proto"},
		},
		{
			files: []string{"a.proto", "b.proto", "c.proto"},
			contents: []string{
				"import 'b.proto';", // a.proto
				"import 'c.proto';", // b.proto
				"import 'a.proto';", // c.proto
			},
			errors: []string{"cycle found: a.proto -> b.proto -> c.proto -> a.proto"},
		},
		{
			files:    []string{"a.proto"},
			contents: []string{"import 'a.proto';"},
			errors:   []string{"cycle found: a.proto -> a.proto"},
		},
	}

	for _, test := range tests {
		units := createUnits(t, test.files, test.contents)

		l := linker.New(units)
		errs := l.Link()

		for i, err := range errs {
			if err.Error() != test.errors[i] {
				t.Fatalf("expected %q, got %q", test.errors[i], err.Error())
			}
		}
	}
}
