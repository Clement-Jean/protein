package parser

import (
	"testing"

	"github.com/Clement-Jean/protein/ast"
	"github.com/Clement-Jean/protein/parser/internal"
	"github.com/Clement-Jean/protein/token"
)

func checkImport(t *testing.T, got ast.Import, expected ast.Import) {
	checkIDs(t, got.ID, expected.ID)
	checkIDs(t, got.Value.ID, expected.Value.ID)

	if got.IsPublic != expected.IsPublic {
		t.Fatalf("expected public import %t, got %t", expected.IsPublic, got.IsPublic)
	}

	if got.IsWeak != expected.IsWeak {
		t.Fatalf("expected weak import %t, got %t", expected.IsWeak, got.IsWeak)
	}
}

func TestImport(t *testing.T) {
	tests := []TestCase[ast.Import]{
		{
			name:        internal.CaseName("import", true),
			expectedObj: &ast.Import{ID: 4, Value: ast.String{ID: 1}},

			content: "import 'google/protobuf/empty.proto';",
			indices: "a-----bc----------------------------de",
			locs:    [][2]rune{{'a', 'b'}, {'c', 'd'}, {'d', 'e'}},
			kinds: []token.Kind{
				token.KindIdentifier,
				token.KindStr,
				token.KindSemicolon,
				token.KindEOF,
			},
		},
		{
			name:        internal.CaseName("import", true, "public"),
			expectedObj: &ast.Import{ID: 5, Value: ast.String{ID: 2}, IsPublic: true},

			content: "import public 'google/protobuf/empty.proto';",
			indices: "a-----bc-----de----------------------------fg",
			locs:    [][2]rune{{'a', 'b'}, {'c', 'd'}, {'e', 'f'}, {'f', 'g'}},
			kinds: []token.Kind{
				token.KindIdentifier, // import
				token.KindIdentifier, // public
				token.KindStr,
				token.KindSemicolon,
				token.KindEOF,
			},
		},
		{
			name:        internal.CaseName("import", true, "weak"),
			expectedObj: &ast.Import{ID: 5, Value: ast.String{ID: 2}, IsWeak: true},

			content: "import weak 'google/protobuf/empty.proto';",
			indices: "a-----bc---de----------------------------fg",
			locs:    [][2]rune{{'a', 'b'}, {'c', 'd'}, {'e', 'f'}, {'f', 'g'}},
			kinds: []token.Kind{
				token.KindIdentifier, // import
				token.KindIdentifier, // weak
				token.KindStr,
				token.KindSemicolon,
				token.KindEOF,
			},
		},
		{
			name: internal.CaseName("import", false, "expect_string"),
			expectedErrs: []error{
				gotUnexpected(&token.Token{ID: 1, Kind: token.KindInt}, token.KindStr),
			},

			content: "import 2;",
			indices: "a-----bcde",
			locs:    [][2]rune{{'a', 'b'}, {'c', 'd'}, {'d', 'e'}},
			kinds: []token.Kind{
				token.KindIdentifier,
				token.KindInt,
				token.KindSemicolon,
				token.KindEOF,
			},
		},
		{
			name: internal.CaseName("import", false, "expect_semicolon"),
			expectedErrs: []error{
				gotUnexpected(&token.Token{ID: 2, Kind: token.KindEOF}, token.KindSemicolon),
			},

			content: "import 'empty.proto'",
			indices: "a-----bc------------d",
			locs:    [][2]rune{{'a', 'b'}, {'c', 'd'}},
			kinds: []token.Kind{
				token.KindIdentifier,
				token.KindStr,
				token.KindEOF,
			},
		},
		{
			name: internal.CaseName("import", false, "expect_public_or_weak"),
			expectedErrs: []error{
				gotUnexpected(&token.Token{ID: 1, Kind: token.KindIdentifier}, token.KindPublic, token.KindWeak),
			},

			content: "import wrong 'empty.proto';",
			indices: "a-----bc----de------------fg",
			locs:    [][2]rune{{'a', 'b'}, {'c', 'd'}, {'e', 'f'}, {'f', 'g'}},
			kinds: []token.Kind{
				token.KindIdentifier, // import
				token.KindIdentifier, // wrong
				token.KindStr,
				token.KindSemicolon,
			},
		},
	}

	wrap := func(p *impl) (ast.Import, []error) {
		imp, err := p.parseImport()
		return imp, internal.EmptyErrorSliceIfNil(err)
	}
	runTestCases(t, tests, checkImport, wrap)
}
