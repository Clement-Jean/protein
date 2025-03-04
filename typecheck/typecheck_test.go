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
)

type pair[K any, V any] struct {
	key K
	val V
}

type testFile = pair[string, string]

func createUnits(t *testing.T, contents []testFile) []*typecheck.Unit {
	t.Helper()

	var units []*typecheck.Unit

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

		units = append(units, &typecheck.Unit{
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
	tests := []typecheckTestCase{
		{
			name: "not defined",
			contents: []testFile{
				{"a.proto", "message A { C c = 1; }"},
			},
			errors: []error{
				&typecheck.TypeNotDefinedError{Name: "C"},
				&typecheck.TypeUnusedWarning{Name: ".A"},
			},
		},
		{
			name: "redefined",
			contents: []testFile{
				{"a.proto", "message A {} message A {}"},
			},
			errors: []error{
				&typecheck.TypeRedefinedError{Name: ".A"},
			},
		},
		{
			name: "redefined across files",
			contents: []testFile{
				{"a.proto", "message A {}"},
				{"b.proto", "message A {}"},
			},
			errors: []error{
				&typecheck.TypeRedefinedError{Name: ".A"},
			},
		},
		{
			name: "import in all includePaths",
			contents: []testFile{
				{"a.proto", "import 'b.proto'; message A { C c = 1; }"},
				{"b.proto", "message C {}"},
				{"test/b.proto", "message B {}"},
			},
			errors: []error{
				&typecheck.TypeUnusedWarning{Name: ".A"},
				&typecheck.TypeUnusedWarning{Name: ".B"},
			},
		},
		{
			name: "import in all includePaths error",
			contents: []testFile{
				{"a.proto", "import 'b.proto'; message A { C c = 1; }"},
				{"b.proto", "message B {}"},
				{"test/b.proto", "message C {}"},
			},
			errors: []error{
				&typecheck.TypeNotImportedError{Name: "C", DefFile: "test/b.proto", RefFile: "a.proto"},
				&typecheck.TypeUnusedWarning{Name: ".A"},
				&typecheck.TypeUnusedWarning{Name: ".B"},
				&typecheck.TypeUnusedWarning{Name: ".C"},
			},
		},
		{
			name: "unknown import in all includePaths",
			contents: []testFile{
				{"a.proto", "import 'b.proto'; message A { C c = 1; }"},
			},
			unknown: []testFile{
				{"b.proto", "message C {}"},
				{"test/b.proto", "message B {}"}, // never parsed
			},
			errors: []error{
				&typecheck.TypeUnusedWarning{Name: ".A"},
			},
		},
		{
			name: "unknown import in all includePaths error",
			contents: []testFile{
				{"a.proto", "import 'b.proto'; message A { C c = 1; }"},
			},
			unknown: []testFile{
				{"b.proto", "message B {}"},
				{"test/b.proto", "message C {}"}, // never parsed
			},
			errors: []error{
				&typecheck.TypeNotDefinedError{Name: "C"},
				&typecheck.TypeUnusedWarning{Name: ".A"},
				&typecheck.TypeUnusedWarning{Name: ".B"},
			},
		},
		{
			name: "use nested type",
			contents: []testFile{
				{"a.proto", "message A { message B {} B b = 1; }"},
			},
			errors: []error{
				&typecheck.TypeUnusedWarning{Name: ".A"},
			},
		},
		{
			name: "map value",
			contents: []testFile{
				{"a.proto", "message A { map<int32, B> b = 1; } message B {}"},
			},
			errors: []error{
				&typecheck.TypeUnusedWarning{Name: ".A"},
			},
		},
		{
			name: "map value across files",
			contents: []testFile{
				{"a.proto", "import 'b.proto'; message A { map<int32, B> b = 1; }"},
				{"b.proto", "message B {}"},
			},
			errors: []error{
				&typecheck.TypeUnusedWarning{Name: ".A"},
			},
		},
		{
			name: "unknown import map value across files",
			contents: []testFile{
				{"a.proto", "import 'b.proto'; message A { map<int32, B> b = 1; }"},
			},
			unknown: []testFile{
				{"b.proto", "message B {}"},
			},
			errors: []error{
				&typecheck.TypeUnusedWarning{Name: ".A"},
			},
		},
		{
			name: "package map value across files",
			contents: []testFile{
				{"a.proto", "import 'b.proto'; message A { map<int32, google.protobuf.B> b = 1; }"},
				{"b.proto", "package google.protobuf; message B {}"},
			},
			errors: []error{
				&typecheck.TypeUnusedWarning{Name: ".A"},
			},
		},
		{
			name: "unknown import package map value across files",
			contents: []testFile{
				{"a.proto", "import 'b.proto'; message A { map<int32, google.protobuf.B> b = 1; }"},
			},
			unknown: []testFile{
				{"b.proto", "package google.protobuf; message B {}"},
			},
			errors: []error{
				&typecheck.TypeUnusedWarning{Name: ".A"},
			},
		},
		{
			name: "unnamed",
			contents: []testFile{
				{"a.proto", "package a.b; message A { message B {} map<int32, B> b = 1; }"},
			},
			errors: []error{
				&typecheck.TypeUnusedWarning{Name: ".a.b.A"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			units := createUnits(t, test.contents)

			t.Run(test.name, func(t *testing.T) {
				l := typecheck.New(
					units,
					typecheck.WithIncludePaths(includePaths...),
					typecheck.WithSourceCreator(fakeSourceCreator(test.contents, test.unknown)),
					typecheck.WithFileCheck(fakeFileCheck(test.contents, test.unknown)),
				)
				errs := l.Check()

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
