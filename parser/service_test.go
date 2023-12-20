package parser

import (
	"testing"

	"github.com/Clement-Jean/protein/ast"
	"github.com/Clement-Jean/protein/parser/internal"
	"github.com/Clement-Jean/protein/token"
)

func checkRpc(t *testing.T, got, expected ast.Rpc) {
	checkIDs(t, got.ID, expected.ID)
	checkIDs(t, got.Name.ID, expected.Name.ID)
	checkIdentifier(t, got.InputType, expected.InputType)
	checkIdentifier(t, got.OutputType, expected.OutputType)

	if got.IsServerStream != expected.IsServerStream {
		t.Fatalf("expected to be server stream (%t), got %t", expected.IsServerStream, got.IsServerStream)
	}

	if got.IsClientStream != expected.IsClientStream {
		t.Fatalf("expected to be client stream (%t), got %t", expected.IsClientStream, got.IsClientStream)
	}

	if len(got.Options) != len(expected.Options) {
		t.Fatalf("expected %d options, got %d", len(expected.Options), len(got.Options))
	}
	for i, opt := range got.Options {
		checkOption(t, opt, expected.Options[i])
	}
}

func checkService(t *testing.T, got, expected ast.Service) {
	checkIDs(t, got.ID, expected.ID)
	checkIDs(t, got.Name.ID, expected.Name.ID)

	if len(got.Options) != len(expected.Options) {
		t.Fatalf("expected %d options, got %d", len(expected.Options), len(got.Options))
	}
	for i, opt := range got.Options {
		checkOption(t, opt, expected.Options[i])
	}

	if len(got.Rpcs) != len(expected.Rpcs) {
		t.Fatalf("expected %d rpc, got %d", len(expected.Rpcs), len(got.Rpcs))
	}
	for i, opt := range got.Rpcs {
		checkRpc(t, opt, expected.Rpcs[i])
	}
}

func TestParseService(t *testing.T) {
	tests := []TestCase[ast.Service]{
		{
			name:        internal.CaseName("service", true),
			expectedObj: ast.Service{ID: 5, Name: ast.Identifier{ID: 1}},

			content: "service Test {}",
			indices: "a------bc---defg",
			locs:    [][2]rune{{'a', 'b'}, {'c', 'd'}, {'e', 'f'}, {'f', 'g'}},
			kinds: []token.Kind{
				token.KindIdentifier, // service
				token.KindIdentifier, // Test
				token.KindLeftBrace,
				token.KindRightBrace,
			},
		},
		{
			name:        internal.CaseName("service", true, "empty_statement"),
			expectedObj: ast.Service{ID: 6, Name: ast.Identifier{ID: 1}},

			content: "service Test {;}",
			indices: "a------bc---defgh",
			locs:    [][2]rune{{'a', 'b'}, {'c', 'd'}, {'e', 'f'}, {'f', 'g'}, {'g', 'h'}},
			kinds: []token.Kind{
				token.KindIdentifier, // service
				token.KindIdentifier, // Test
				token.KindLeftBrace,
				token.KindSemicolon,
				token.KindRightBrace,
			},
		},
		{
			name: internal.CaseName("service", true, "option"),
			expectedObj: ast.Service{
				ID: 11, Name: ast.Identifier{ID: 1},
				Options: []ast.Option{{
					ID: 10, Name: ast.Identifier{ID: 4}, Value: &ast.Boolean{ID: 6},
				}},
			},

			content: "service Test { option deprecated = true; }",
			indices: "a------bc---defg-----hi---------jklm---nopq",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'e', 'f'}, {'g', 'h'},
				{'i', 'j'}, {'k', 'l'}, {'m', 'n'}, {'n', 'o'},
				{'p', 'q'},
			},
			kinds: []token.Kind{
				token.KindIdentifier, // service
				token.KindIdentifier, // Test
				token.KindLeftBrace,
				token.KindIdentifier,
				token.KindIdentifier, // deprecated
				token.KindEqual,
				token.KindIdentifier, // true
				token.KindSemicolon,
				token.KindRightBrace,
			},
		},
		{
			name: internal.CaseName("service", true, "rpc"),
			expectedObj: ast.Service{
				ID: 16, Name: ast.Identifier{ID: 1},
				Rpcs: []ast.Rpc{
					{ID: 15, Name: ast.Identifier{ID: 4}, InputType: ast.Identifier{ID: 6}, OutputType: ast.Identifier{ID: 10}},
				},
			},

			content: "service Test { rpc T (Test) returns (Test); }",
			indices: "a------bc---defg--hijkl---mno------pqr---stuvw",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'e', 'f'}, {'g', 'h'},
				{'i', 'j'}, {'k', 'l'}, {'l', 'm'}, {'m', 'n'},
				{'o', 'p'}, {'q', 'r'}, {'r', 's'}, {'s', 't'},
				{'t', 'u'}, {'v', 'w'},
			},
			kinds: []token.Kind{
				token.KindIdentifier, // service
				token.KindIdentifier, // Test
				token.KindLeftBrace,
				token.KindIdentifier, // rpc
				token.KindIdentifier, // T
				token.KindLeftParen,
				token.KindIdentifier, // Test
				token.KindRightParen,
				token.KindIdentifier, // returns
				token.KindLeftParen,
				token.KindIdentifier, // Test
				token.KindRightParen,
				token.KindSemicolon,
				token.KindRightBrace,
			},
		},
		{
			name: internal.CaseName("service", false, "expected_identifier"),
			expectedErrs: []error{
				gotUnexpected(&token.Token{ID: 1, Kind: token.KindLeftBrace}, token.KindIdentifier),
			},

			content: "service {}",
			indices: "a------bcde",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'d', 'e'},
			},
			kinds: []token.Kind{
				token.KindIdentifier, // service
				token.KindLeftBrace,
				token.KindRightBrace,
			},
		},
		{
			name: internal.CaseName("service", false, "expected_left_brace"),
			expectedErrs: []error{
				gotUnexpected(&token.Token{ID: 2, Kind: token.KindLeftSquare}, token.KindLeftBrace),
			},

			content: "service Test [}",
			indices: "a------bc---defg",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'d', 'e'}, {'f', 'g'},
			},
			kinds: []token.Kind{
				token.KindIdentifier, // service
				token.KindIdentifier, // Test
				token.KindLeftSquare,
				token.KindRightBrace,
			},
		},
		{
			name: internal.CaseName("service", false, "expected_right_brace"),
			expectedErrs: []error{
				gotUnexpected(&token.Token{ID: 3, Kind: token.KindEOF}, token.KindRightBrace),
			},

			content: "service Test {",
			indices: "a------bc---def",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'e', 'f'},
			},
			kinds: []token.Kind{
				token.KindIdentifier, // service
				token.KindIdentifier, // Test
				token.KindLeftBrace,
			},
		},
		{
			name: internal.CaseName("service", false, "unexpected_int"),
			expectedErrs: []error{
				gotUnexpected(&token.Token{ID: 3, Kind: token.KindInt}, token.KindOption, token.KindRpc, token.KindRightBrace),
			},

			content: "service Test { 2 }",
			indices: "a------bc---defghij",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'e', 'f'}, {'g', 'h'},
				{'i', 'j'},
			},
			kinds: []token.Kind{
				token.KindIdentifier, // service
				token.KindIdentifier, // Test
				token.KindLeftBrace,
				token.KindInt,
				token.KindRightBrace,
			},
		},
	}

	runTestCases(t, tests, checkService, (*impl).parseService)
}

func TestParseRpc(t *testing.T) {
	tests := []TestCase[ast.Rpc]{
		{
			name: internal.CaseName("rpc", true),
			expectedObj: ast.Rpc{
				ID: 11, Name: ast.Identifier{ID: 1}, InputType: ast.Identifier{ID: 3}, OutputType: ast.Identifier{ID: 7},
			},

			content: "rpc T (Test) returns (Test);",
			indices: "a--bcdef---ghi------jkl---mno",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'e', 'f'}, {'f', 'g'},
				{'g', 'h'}, {'i', 'j'}, {'k', 'l'}, {'l', 'm'},
				{'m', 'n'}, {'n', 'o'},
			},
			kinds: []token.Kind{
				token.KindIdentifier, // rpc
				token.KindIdentifier, // T
				token.KindLeftParen,
				token.KindIdentifier, // Test
				token.KindRightParen,
				token.KindIdentifier, // returns
				token.KindLeftParen,
				token.KindIdentifier, // Test
				token.KindRightParen,
				token.KindSemicolon,
			},
		},
		{
			name: internal.CaseName("rpc", true, "option"),
			expectedObj: ast.Rpc{
				ID: 18, Name: ast.Identifier{ID: 1}, InputType: ast.Identifier{ID: 3}, OutputType: ast.Identifier{ID: 7},
				Options: []ast.Option{
					{ID: 17, Name: ast.Identifier{ID: 11}, Value: &ast.Boolean{ID: 13}},
				},
			},

			content: "rpc T (Test) returns (Test) { option deprecated = true; }",
			indices: "a--bcdef---ghi------jkl---mnopq-----rs---------tuvw---xyz1",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'e', 'f'}, {'f', 'g'},
				{'g', 'h'}, {'i', 'j'}, {'k', 'l'}, {'l', 'm'},
				{'m', 'n'}, {'o', 'p'}, {'q', 'r'}, {'s', 't'},
				{'u', 'v'}, {'w', 'x'}, {'x', 'y'}, {'z', '1'},
			},
			kinds: []token.Kind{
				token.KindIdentifier, // rpc
				token.KindIdentifier, // T
				token.KindLeftParen,
				token.KindIdentifier, // Test
				token.KindRightParen,
				token.KindIdentifier, // returns
				token.KindLeftParen,
				token.KindIdentifier, // Test
				token.KindRightParen,
				token.KindLeftBrace,
				token.KindIdentifier, // option
				token.KindIdentifier, // deprecated
				token.KindEqual,
				token.KindIdentifier, // true
				token.KindSemicolon,
				token.KindRightBrace,
			},
		},
		{
			name: internal.CaseName("rpc", true, "empty_statement_in_block"),
			expectedObj: ast.Rpc{
				ID: 13, Name: ast.Identifier{ID: 1}, InputType: ast.Identifier{ID: 3}, OutputType: ast.Identifier{ID: 7},
			},

			content: "rpc T (Test) returns (Test) { ; }",
			indices: "a--bcdef---ghi------jkl---mnopqrst",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'e', 'f'}, {'f', 'g'},
				{'g', 'h'}, {'i', 'j'}, {'k', 'l'}, {'l', 'm'},
				{'m', 'n'}, {'o', 'p'}, {'q', 'r'}, {'s', 't'},
			},
			kinds: []token.Kind{
				token.KindIdentifier, // rpc
				token.KindIdentifier, // T
				token.KindLeftParen,
				token.KindIdentifier, // Test
				token.KindRightParen,
				token.KindIdentifier, // returns
				token.KindLeftParen,
				token.KindIdentifier, // Test
				token.KindRightParen,
				token.KindLeftBrace,
				token.KindSemicolon,
				token.KindRightBrace,
			},
		},
		{
			name: internal.CaseName("rpc", true, "stream"),
			expectedObj: ast.Rpc{
				ID: 13, Name: ast.Identifier{ID: 1},
				InputType: ast.Identifier{ID: 4}, OutputType: ast.Identifier{ID: 9},
				IsServerStream: true, IsClientStream: true,
			},

			content: "rpc T (stream Test) returns (stream Test);",
			indices: "a--bcdef-----gh---ijk------lmn-----op---qrs",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'e', 'f'}, {'f', 'g'},
				{'h', 'i'}, {'i', 'j'}, {'k', 'l'}, {'m', 'n'},
				{'n', 'o'}, {'p', 'q'}, {'q', 'r'}, {'r', 's'},
			},
			kinds: []token.Kind{
				token.KindIdentifier, // rpc
				token.KindIdentifier, // T
				token.KindLeftParen,
				token.KindIdentifier, // stream
				token.KindIdentifier, // Test
				token.KindRightParen,
				token.KindIdentifier, // returns
				token.KindLeftParen,
				token.KindIdentifier, // stream
				token.KindIdentifier, // Test
				token.KindRightParen,
				token.KindSemicolon,
			},
		},
		{
			name: internal.CaseName("rpc", false, "unexpected_int"),
			expectedErrs: []error{
				gotUnexpected(&token.Token{ID: 10, Kind: token.KindInt}, token.KindOption, token.KindRightBrace),
			},

			content: "rpc T (Test) returns (Test) { 2 }",
			indices: "a--bcdef---ghi------jkl---mnopqrst",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'e', 'f'}, {'f', 'g'},
				{'g', 'h'}, {'i', 'j'}, {'k', 'l'}, {'l', 'm'},
				{'m', 'n'}, {'o', 'p'}, {'q', 'r'}, {'s', 't'},
			},
			kinds: []token.Kind{
				token.KindIdentifier, // rpc
				token.KindIdentifier, // T
				token.KindLeftParen,
				token.KindIdentifier, // Test
				token.KindRightParen,
				token.KindIdentifier, // returns
				token.KindLeftParen,
				token.KindIdentifier, // Test
				token.KindRightParen,
				token.KindLeftBrace,
				token.KindInt,
				token.KindRightBrace,
			},
		},
		{
			name: internal.CaseName("rpc", false, "expected_identifier"),
			expectedErrs: []error{
				gotUnexpected(&token.Token{ID: 1, Kind: token.KindLeftParen}, token.KindIdentifier),
			},

			content: "rpc (Test) returns (Test);",
			indices: "a--bcd---efg------hij---klm",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'d', 'e'}, {'e', 'f'},
				{'g', 'h'}, {'i', 'j'}, {'j', 'k'}, {'k', 'l'},
				{'l', 'm'},
			},
			kinds: []token.Kind{
				token.KindIdentifier, // rpc
				token.KindLeftParen,
				token.KindIdentifier, // Test
				token.KindRightParen,
				token.KindIdentifier, // returns
				token.KindLeftParen,
				token.KindIdentifier, // Test
				token.KindRightParen,
				token.KindSemicolon,
			},
		},
		{
			name: internal.CaseName("rpc", false, "expected_in_left_paren"),
			expectedErrs: []error{
				gotUnexpected(&token.Token{ID: 2, Kind: token.KindLeftSquare}, token.KindLeftParen),
			},

			content: "rpc T [Test) returns (Test);",
			indices: "a--bcdef---ghi------jkl---mno",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'e', 'f'}, {'f', 'g'},
				{'g', 'h'}, {'i', 'j'}, {'k', 'l'}, {'l', 'm'},
				{'m', 'n'}, {'n', 'o'},
			},
			kinds: []token.Kind{
				token.KindIdentifier, // rpc
				token.KindIdentifier, // T
				token.KindLeftSquare,
				token.KindIdentifier, // Test
				token.KindRightParen,
				token.KindIdentifier, // returns
				token.KindLeftParen,
				token.KindIdentifier, // Test
				token.KindRightParen,
				token.KindSemicolon,
			},
		},
		{
			name: internal.CaseName("rpc", false, "expected_in_right_paren"),
			expectedErrs: []error{
				gotUnexpected(&token.Token{ID: 4, Kind: token.KindRightSquare}, token.KindRightParen),
			},

			content: "rpc T (Test] returns (Test);",
			indices: "a--bcdef---ghi------jkl---mno",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'e', 'f'}, {'f', 'g'},
				{'g', 'h'}, {'i', 'j'}, {'k', 'l'}, {'l', 'm'},
				{'m', 'n'}, {'n', 'o'},
			},
			kinds: []token.Kind{
				token.KindIdentifier, // rpc
				token.KindIdentifier, // T
				token.KindLeftParen,
				token.KindIdentifier, // Test
				token.KindRightSquare,
				token.KindIdentifier, // returns
				token.KindLeftParen,
				token.KindIdentifier, // Test
				token.KindRightParen,
				token.KindSemicolon,
			},
		},
		{
			name: internal.CaseName("rpc", false, "expected_in_identifier"),
			expectedErrs: []error{
				gotUnexpected(&token.Token{ID: 3, Kind: token.KindStr}, token.KindIdentifier),
			},

			content: "rpc T ('Test') returns (Test);",
			indices: "a--bcdef-----ghi------jkl---mno",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'e', 'f'}, {'f', 'g'},
				{'g', 'h'}, {'i', 'j'}, {'k', 'l'}, {'l', 'm'},
				{'m', 'n'}, {'n', 'o'},
			},
			kinds: []token.Kind{
				token.KindIdentifier, // rpc
				token.KindIdentifier, // T
				token.KindLeftParen,
				token.KindStr,
				token.KindRightParen,
				token.KindIdentifier, // returns
				token.KindLeftParen,
				token.KindIdentifier, // Test
				token.KindRightParen,
				token.KindSemicolon,
			},
		},
		{
			name: internal.CaseName("rpc", false, "expected_returns"),
			expectedErrs: []error{
				gotUnexpected(&token.Token{ID: 5, Kind: token.KindIdentifier}, token.KindReturns),
			},

			content: "rpc T (Test) return (Test);",
			indices: "a--bcdef---ghi-----jkl---mno",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'e', 'f'}, {'f', 'g'},
				{'g', 'h'}, {'i', 'j'}, {'k', 'l'}, {'l', 'm'},
				{'m', 'n'}, {'n', 'o'},
			},
			kinds: []token.Kind{
				token.KindIdentifier, // rpc
				token.KindIdentifier, // T
				token.KindLeftParen,
				token.KindIdentifier, // Test
				token.KindRightParen,
				token.KindIdentifier, // return
				token.KindLeftParen,
				token.KindIdentifier, // Test
				token.KindRightParen,
				token.KindSemicolon,
			},
		},
		{
			name: internal.CaseName("rpc", false, "expected_returns_identifier"),
			expectedErrs: []error{
				gotUnexpected(&token.Token{ID: 5, Kind: token.KindStr}, token.KindReturns),
			},

			content: "rpc T (Test) 'return' (Test);",
			indices: "a--bcdef---ghi-------jkl---mno",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'e', 'f'}, {'f', 'g'},
				{'g', 'h'}, {'i', 'j'}, {'k', 'l'}, {'l', 'm'},
				{'m', 'n'}, {'n', 'o'},
			},
			kinds: []token.Kind{
				token.KindIdentifier, // rpc
				token.KindIdentifier, // T
				token.KindLeftParen,
				token.KindIdentifier, // Test
				token.KindRightParen,
				token.KindStr,
				token.KindLeftParen,
				token.KindIdentifier, // Test
				token.KindRightParen,
				token.KindSemicolon,
			},
		},
		{
			name: internal.CaseName("rpc", false, "expected_out_left_paren"),
			expectedErrs: []error{
				gotUnexpected(&token.Token{ID: 6, Kind: token.KindLeftSquare}, token.KindLeftParen),
			},

			content: "rpc T (Test) returns [Test);",
			indices: "a--bcdef---ghi------jkl---mno",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'e', 'f'}, {'f', 'g'},
				{'g', 'h'}, {'i', 'j'}, {'k', 'l'}, {'l', 'm'},
				{'m', 'n'}, {'n', 'o'},
			},
			kinds: []token.Kind{
				token.KindIdentifier, // rpc
				token.KindIdentifier, // T
				token.KindLeftParen,
				token.KindIdentifier, // Test
				token.KindRightParen,
				token.KindIdentifier, // returns
				token.KindLeftSquare,
				token.KindIdentifier, // Test
				token.KindRightParen,
				token.KindSemicolon,
			},
		},
		{
			name: internal.CaseName("rpc", false, "expected_out_right_paren"),
			expectedErrs: []error{
				gotUnexpected(&token.Token{ID: 8, Kind: token.KindRightSquare}, token.KindRightParen),
			},

			content: "rpc T (Test) returns (Test];",
			indices: "a--bcdef---ghi------jkl---mno",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'e', 'f'}, {'f', 'g'},
				{'g', 'h'}, {'i', 'j'}, {'k', 'l'}, {'l', 'm'},
				{'m', 'n'}, {'n', 'o'},
			},
			kinds: []token.Kind{
				token.KindIdentifier, // rpc
				token.KindIdentifier, // T
				token.KindLeftParen,
				token.KindIdentifier, // Test
				token.KindRightParen,
				token.KindIdentifier, // returns
				token.KindLeftParen,
				token.KindIdentifier, // Test
				token.KindRightSquare,
				token.KindSemicolon,
			},
		},
		{
			name: internal.CaseName("rpc", false, "expected_out_identifier"),
			expectedErrs: []error{
				gotUnexpected(&token.Token{ID: 7, Kind: token.KindStr}, token.KindIdentifier),
			},

			content: "rpc T (Test) returns ('Test');",
			indices: "a--bcdef---ghi------jkl-----mno",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'e', 'f'}, {'f', 'g'},
				{'g', 'h'}, {'i', 'j'}, {'k', 'l'}, {'l', 'm'},
				{'m', 'n'}, {'n', 'o'},
			},
			kinds: []token.Kind{
				token.KindIdentifier, // rpc
				token.KindIdentifier, // T
				token.KindLeftParen,
				token.KindIdentifier, // Test
				token.KindRightParen,
				token.KindIdentifier, // returns
				token.KindLeftParen,
				token.KindStr,
				token.KindRightParen,
				token.KindSemicolon,
			},
		},
		{
			name: internal.CaseName("rpc", false, "expected_identifier_after_stream"),
			expectedErrs: []error{
				gotUnexpected(&token.Token{ID: 4, Kind: token.KindStr}, token.KindIdentifier),
			},

			content: "rpc T (stream 'Test') returns ('Test');",
			indices: "a--bcdef-----gh-----ijk------lmn-----opq",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'e', 'f'}, {'f', 'g'},
				{'h', 'i'}, {'i', 'j'}, {'k', 'l'}, {'m', 'n'},
				{'n', 'o'}, {'o', 'p'}, {'p', 'q'},
			},
			kinds: []token.Kind{
				token.KindIdentifier, // rpc
				token.KindIdentifier, // T
				token.KindLeftParen,
				token.KindIdentifier, // stream
				token.KindStr,
				token.KindRightParen,
				token.KindIdentifier, // returns
				token.KindLeftParen,
				token.KindIdentifier,
				token.KindRightParen,
				token.KindSemicolon,
			},
		},
		{
			name: internal.CaseName("rpc", true, "expected_right_brace"),
			expectedErrs: []error{
				gotUnexpected(&token.Token{ID: 10, Kind: token.KindEOF}, token.KindRightBrace),
			},

			content: "rpc T (Test) returns (Test) {",
			indices: "a--bcdef---ghi------jkl---mnop",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'e', 'f'}, {'f', 'g'},
				{'g', 'h'}, {'i', 'j'}, {'k', 'l'}, {'l', 'm'},
				{'m', 'n'}, {'o', 'p'},
			},
			kinds: []token.Kind{
				token.KindIdentifier, // rpc
				token.KindIdentifier, // T
				token.KindLeftParen,
				token.KindIdentifier, // Test
				token.KindRightParen,
				token.KindIdentifier, // returns
				token.KindLeftParen,
				token.KindIdentifier, // Test
				token.KindRightParen,
				token.KindLeftBrace,
			},
		},
		{
			name: internal.CaseName("rpc", true, "expected_semicolon"),
			expectedErrs: []error{
				gotUnexpected(&token.Token{ID: 9, Kind: token.KindEOF}, token.KindSemicolon),
			},

			content: "rpc T (Test) returns (Test)",
			indices: "a--bcdef---ghi------jkl---mn",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'e', 'f'}, {'f', 'g'},
				{'g', 'h'}, {'i', 'j'}, {'k', 'l'}, {'l', 'm'},
				{'m', 'n'},
			},
			kinds: []token.Kind{
				token.KindIdentifier, // rpc
				token.KindIdentifier, // T
				token.KindLeftParen,
				token.KindIdentifier, // Test
				token.KindRightParen,
				token.KindIdentifier, // returns
				token.KindLeftParen,
				token.KindIdentifier, // Test
				token.KindRightParen,
			},
		},
	}

	runTestCases(t, tests, checkRpc, (*impl).parseRpc)
}
