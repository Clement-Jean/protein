package typecheck_test

import (
	"github.com/Clement-Jean/protein/lexer"
	"github.com/Clement-Jean/protein/parser"
	"github.com/Clement-Jean/protein/typecheck"
)

var importTests = []typecheckTestCase{
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
	},
	{
		name: "weak import",
		contents: []testFile{
			{"a.proto", "import weak 'b.proto';"},
			{"b.proto", ""},
		},
		errorLevel: typecheck.ErrorLevelWarning,
		errors: []error{
			&typecheck.WeakImportNoEffectWarning{
				File: "a.proto",
				Line: 1,
				Col:  13,
			},
		},
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
				DefFile: "d.proto",
				RefFile: "a.proto",
				Line:    1,
				Col:     31,
				Name:    "a.b.c.d.D",
			},
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
	{
		name: "import in all includePaths",
		contents: []testFile{
			{"a.proto", "import 'b.proto'; message A { C c = 1; }"},
			{"b.proto", "message C {}"},
			{"test/b.proto", "message B {}"},
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
			&typecheck.TypeNotImportedError{
				Name:    "C",
				DefFile: "test/b.proto",
				RefFile: "a.proto",
				Line:    1,
				Col:     31,
			},
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
			&typecheck.TypeNotDefinedError{
				File: "a.proto",
				Line: 1,
				Col:  31,
				Name: "C",
			},
		},
	},
}
