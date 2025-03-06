package typecheck_test

import (
	"testing"

	"github.com/Clement-Jean/protein/lexer"
	"github.com/Clement-Jean/protein/parser"
	"github.com/Clement-Jean/protein/typecheck"
)

type importTestCase struct {
	name     string
	contents []testFile
	unknown  []testFile
	errors   []error
}

func TestImports(t *testing.T) {
	// TODO: absolute path for inputs

	includePaths := []string{"", "test"}
	tests := []importTestCase{
		{
			name: "trivial imports",
			contents: []testFile{
				{"a.proto", "import 'b.proto';"},
				{"b.proto", "import 'c.proto';"},
				{"c.proto", ""},
			},
		},
		{
			name: "trivial imports relative path",
			contents: []testFile{
				{"test/a.proto", "import 'b.proto';"},
				{"test/b.proto", "import 'c.proto';"},
				{"test/c.proto", ""},
			},
		},
		{
			name: "cycle length 1",
			contents: []testFile{
				{"a.proto", "import 'b.proto';"},
				{"b.proto", "import 'a.proto';"},
			},
			errors: []error{&typecheck.ImportCycleError{Files: []string{"a.proto", "b.proto"}}},
		},
		{
			name: "cycle length 2",
			contents: []testFile{
				{"a.proto", "import 'b.proto';"},
				{"b.proto", "import 'c.proto';"},
				{"c.proto", "import 'a.proto';"},
			},
			errors: []error{&typecheck.ImportCycleError{Files: []string{"a.proto", "b.proto", "c.proto"}}},
		},
		{
			name: "cycle length 0",
			contents: []testFile{
				{"a.proto", "import 'a.proto';"},
			},
			errors: []error{&typecheck.ImportCycleError{Files: []string{"a.proto"}}},
		},
		{
			name: "trivial unknown import",
			contents: []testFile{
				{"a.proto", "import 'b.proto';"},
			},
			unknown: []testFile{
				{"b.proto", ""},
			},
		},
		{
			name: "cycle length 1 unknown import",
			contents: []testFile{
				{"a.proto", "import 'b.proto';"},
			},
			unknown: []testFile{
				{"b.proto", "import 'a.proto';"},
			},
			errors: []error{&typecheck.ImportCycleError{Files: []string{"a.proto", "b.proto"}}},
		},
		{
			name: "cycle length 2 unknown import",
			contents: []testFile{
				{"a.proto", "import 'b.proto';"},
			},
			unknown: []testFile{
				{"b.proto", "import 'c.proto';"},
				{"c.proto", "import 'a.proto';"},
			},
			errors: []error{&typecheck.ImportCycleError{Files: []string{"a.proto", "b.proto", "c.proto"}}},
		},
		{
			name: "import not found",
			contents: []testFile{
				{"a.proto", "import 'b.proto';"},
			},
			errors: []error{&typecheck.ImportFileNotFoundError{File: "b.proto"}},
		},
		{
			name: "public import",
			contents: []testFile{
				{"a.proto", "import 'b.proto';"},
				{"b.proto", "import 'c.proto'; message B { a.b.c.d.D d = 1; }"},
				{"c.proto", "import public 'd.proto';"},
				{"d.proto", "package a.b.c.d; message D {}"},
			},
			errors: []error{&typecheck.TypeUnusedWarning{Name: ".B"}},
		},
		{
			name: "weak import",
			contents: []testFile{
				{"a.proto", "import weak 'b.proto';"},
				{"b.proto", ""},
			},
			errors: []error{&typecheck.WeakImportNoEffectWarning{}},
		},
		{
			name: "public import not accessible error",
			contents: []testFile{
				{"a.proto", "import 'b.proto'; message A { a.b.c.d.D d = 1; }"},
				{"b.proto", "import 'c.proto';"},
				{"c.proto", "import public 'd.proto';"},
				{"d.proto", "package a.b.c.d; message D {}"},
			},
			errors: []error{
				&typecheck.TypeNotImportedError{
					Name:    "a.b.c.d.D",
					DefFile: "d.proto",
					RefFile: "a.proto",
				},
				&typecheck.TypeUnusedWarning{Name: ".A"},
				&typecheck.TypeUnusedWarning{Name: ".a.b.c.d.D"},
			},
		},
		{
			name: "unknown import file with lexer error",
			contents: []testFile{
				{"a.proto", "import 'b.proto';"},
			},
			unknown: []testFile{
				{"b.proto", "错"},
			},
			errors: []error{
				&lexer.InvalidChar{Character: "错"[0]},
				&lexer.InvalidChar{Character: "错"[1]},
				&lexer.InvalidChar{Character: "错"[2]},
				&typecheck.ImportFileNotFoundError{File: "b.proto"},
			},
		},
		{
			name: "unknown import file with parse error",
			contents: []testFile{
				{"a.proto", "import 'b.proto';"},
			},
			unknown: []testFile{
				{"b.proto", "message {}"},
			},
			errors: []error{
				&parser.ExpectedError{
					Expected: []lexer.TokenKind{lexer.TokenKindIdentifier},
					Got:      lexer.TokenKindLeftBrace,
				},
				&typecheck.ImportFileNotFoundError{File: "b.proto"},
			},
		},
	}

	for _, test := range tests {
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
	}
}
