package typecheck_test

import (
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/Clement-Jean/protein/lexer"
	"github.com/Clement-Jean/protein/parser"
	"github.com/Clement-Jean/protein/source"
	"github.com/Clement-Jean/protein/typecheck"
	"github.com/Clement-Jean/protein/unit"
)

type pair[K any, V any] struct {
	key K
	val V
}

type testFile = pair[string, string]

func createUnits(t *testing.T, contents []testFile) []*unit.Unit {
	t.Helper()

	var units []*unit.Unit

	// sort for being able to use binary search
	slices.SortFunc(contents, func(p, p2 testFile) int {
		return strings.Compare(p.key, p2.key)
	})

	for _, pair := range contents {
		file := pair.key
		content := pair.val
		s, err := source.NewFromReader(strings.NewReader(content))
		if err != nil {
			t.Fatal(err)
		}

		l := lexer.NewFromSource(s)
		tb, errs := l.Lex()
		if len(errs) != 0 {
			t.Fatal(errs)
		}

		p := parser.New(tb)
		pt, errs := p.Parse()
		if len(errs) != 0 {
			t.Fatal(errs)
		}

		units = append(units, &unit.Unit{
			File:   file,
			Buffer: s,
			Toks:   tb,
			Tree:   pt,
		})
	}

	return units
}

func findPath(a string) func(testFile) bool {
	return func(b testFile) bool {
		return b.key == a
	}
}

func fakeSourceCreator(contents, unknown []testFile) typecheck.SourceCreator {
	return func(path string) (*source.Buffer, error) {
		if idx := slices.IndexFunc(contents, findPath(path)); idx != -1 {
			return source.NewFromReader(strings.NewReader(contents[idx].val))
		}

		if idx := slices.IndexFunc(unknown, findPath(path)); idx != -1 {
			return source.NewFromReader(strings.NewReader(unknown[idx].val))
		}

		return nil, os.ErrNotExist
	}
}

func fakeFileCheck(contents, unknown []testFile) typecheck.FileExistsCheck {
	return func(path string) bool {
		return slices.ContainsFunc(contents, findPath(path)) || slices.ContainsFunc(unknown, findPath(path))
	}
}

type typecheckTestCase struct {
	name     string
	contents []testFile
	unknown  []testFile
	errors   []error
}

// TODO: set error level per test!

func TestTypeCheck(t *testing.T) {
	includePaths := []string{"", "test"}
	tests := slices.Concat(
		packageTests,
		importTests,
		messageTests,
		mapTests,
		oneofTests,
		enumTests,
		serviceTests,
	)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			units := createUnits(t, test.contents)

			t.Run(test.name, func(t *testing.T) {
				l := typecheck.New(
					units,
					typecheck.WithIncludePaths(includePaths...),
					typecheck.WithSourceCreator(fakeSourceCreator(test.contents, test.unknown)),
					typecheck.WithFileCheck(fakeFileCheck(test.contents, test.unknown)),
					typecheck.WithErrorLevel(typecheck.ErrorLevelError),
				)
				_, errs := l.Check()

				if len(errs) != len(test.errors) {
					t.Fatalf("expected %d errors, got %d: %v", len(test.errors), len(errs), errs)
				}

				for i, err := range errs {
					if err.Error() != test.errors[i].Error() {
						t.Fatalf("expected %q, got %q", test.errors[i].Error(), err.Error())
					}
				}
			})
		})
	}
}
