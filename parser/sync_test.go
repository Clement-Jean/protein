package parser

import (
	"testing"

	"github.com/Clement-Jean/protein/ast"
	"github.com/Clement-Jean/protein/parser/internal"
	"github.com/Clement-Jean/protein/token"
)

func TestSyncOption(t *testing.T) {
	tests := []TestCase[ast.Option]{
		{
			name: internal.CaseName("sync", true, "text_message"),
			expectedObj: &ast.Option{
				ID:   14,
				Name: ast.Identifier{ID: 1},
				Value: ast.TextMessage{
					ID: 13,
					Fields: []ast.TextField{
						{
							ID:    12,
							Name:  ast.Identifier{ID: 6},
							Value: ast.Identifier{ID: 8},
						},
					},
				},
			},
			expectedErrs: []error{
				gotUnexpected(&token.Token{ID: 4, Kind: token.KindStr}, token.KindIdentifier, token.KindLeftSquare),
			},

			content: "option test = { 'a', b: c };",
			indices: "a-----bc---defghi--jklmnopqrs",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'e', 'f'}, {'g', 'h'},
				{'i', 'j'}, {'j', 'k'}, {'l', 'm'}, {'m', 'n'},
				{'o', 'p'}, {'q', 'r'}, {'r', 's'},
			},
			kinds: []token.Kind{
				token.KindIdentifier, // option
				token.KindIdentifier, // test
				token.KindEqual,
				token.KindLeftBrace,
				token.KindStr,
				token.KindComma,
				token.KindIdentifier, // b
				token.KindColon,
				token.KindIdentifier, // c
				token.KindRightBrace,
				token.KindSemicolon,
			},
		},
		{
			name: internal.CaseName("sync", true, "text_message_with_invalid_constant"),
			expectedObj: &ast.Option{
				ID:   14,
				Name: ast.Identifier{ID: 1},
				Value: ast.TextMessage{
					ID: 13,
					Fields: []ast.TextField{
						{
							ID:    12,
							Name:  ast.Identifier{ID: 6},
							Value: ast.Identifier{ID: 8},
						},
					},
				},
			},
			expectedErrs: []error{
				gotUnexpected(&token.Token{ID: 4, Kind: token.KindColon}, token.KindIdentifier, token.KindLeftSquare),
			},

			content: "option test = { :, b: c };",
			indices: "a-----bc---defghijklmnopqrs",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'e', 'f'}, {'g', 'h'},
				{'i', 'j'}, {'j', 'k'}, {'l', 'm'}, {'m', 'n'},
				{'o', 'p'}, {'q', 'r'}, {'r', 's'},
			},
			kinds: []token.Kind{
				token.KindIdentifier, // option
				token.KindIdentifier, // test
				token.KindEqual,
				token.KindLeftBrace,
				token.KindColon,
				token.KindComma,
				token.KindIdentifier, // b
				token.KindColon,
				token.KindIdentifier, // c
				token.KindRightBrace,
				token.KindSemicolon,
			},
		},
	}

	runTestCases(t, tests, checkOption, (*impl).parseOption)
}

func TestSyncInlineOption(t *testing.T) {
	tests := []TestCase[[]ast.Option]{
		{
			name: internal.CaseName("sync", true, "inline_option"),
			expectedObj: &[]ast.Option{
				{
					ID:    8,
					Name:  ast.Identifier{ID: 3},
					Value: ast.Identifier{ID: 5},
				},
			},
			expectedErrs: []error{
				gotUnexpected(&token.Token{ID: 1, Kind: token.KindStr}, token.KindIdentifier),
			},

			content: "['a', b = c]",
			indices: "ab--cdefghijk",
			locs: [][2]rune{
				{'a', 'b'}, {'b', 'c'}, {'c', 'd'}, {'e', 'f'},
				{'g', 'h'}, {'i', 'j'}, {'j', 'k'},
			},
			kinds: []token.Kind{
				token.KindLeftSquare,
				token.KindStr,
				token.KindComma,
				token.KindIdentifier, // b
				token.KindEqual,
				token.KindIdentifier, // c
				token.KindRightSquare,
			},
		},
	}

	runTestCases(t, tests, checkOptions, (*impl).parseInlineOptions)
}

func TestSyncMessage(t *testing.T) {
	tests := []TestCase[ast.Message]{
		{
			name: internal.CaseName("sync", true, "message"),
			expectedObj: &ast.Message{
				ID:   6,
				Name: ast.Identifier{ID: 1},
			},
			expectedErrs: []error{
				gotUnexpected(&token.Token{ID: 3, Kind: token.KindStr}, token.KindOption, token.KindReserved, token.KindField),
			},

			content: "message Test { 'a' }",
			indices: "a------bc---defg--hij",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'e', 'f'},
				{'g', 'h'}, {'i', 'j'},
			},
			kinds: []token.Kind{
				token.KindIdentifier,
				token.KindIdentifier, // Test
				token.KindLeftBrace,
				token.KindStr,
				token.KindRightBrace,
			},
		},
	}

	wrap := func(p *impl) (ast.Message, []error) { return p.parseMessage(1) }
	runTestCases(t, tests, checkMessage, wrap)
}

func TestSyncEnum(t *testing.T) {
	tests := []TestCase[ast.Enum]{
		{
			name: internal.CaseName("sync", true, "enum"),
			expectedObj: &ast.Enum{
				ID:   6,
				Name: ast.Identifier{ID: 1},
			},
			expectedErrs: []error{
				gotUnexpected(&token.Token{ID: 3, Kind: token.KindStr}, token.KindOption, token.KindReserved, token.KindIdentifier),
			},

			content: "enum Test { 'a' }",
			indices: "a---bc---defg--hij",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'e', 'f'},
				{'g', 'h'}, {'i', 'j'},
			},
			kinds: []token.Kind{
				token.KindIdentifier,
				token.KindIdentifier, // Test
				token.KindLeftBrace,
				token.KindStr,
				token.KindRightBrace,
			},
		},
	}

	runTestCases(t, tests, checkEnum, (*impl).parseEnum)
}

func TestSyncService(t *testing.T) {
	tests := []TestCase[ast.Service]{
		{
			name: internal.CaseName("sync", true, "service"),
			expectedObj: &ast.Service{
				ID:   6,
				Name: ast.Identifier{ID: 1},
			},
			expectedErrs: []error{
				gotUnexpected(&token.Token{ID: 3, Kind: token.KindStr}, token.KindOption, token.KindRpc),
			},

			content: "service Test { 'a' }",
			indices: "a------bc---defg--hij",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'e', 'f'},
				{'g', 'h'}, {'i', 'j'},
			},
			kinds: []token.Kind{
				token.KindIdentifier,
				token.KindIdentifier, // Test
				token.KindLeftBrace,
				token.KindStr,
				token.KindRightBrace,
			},
		},
	}

	runTestCases(t, tests, checkService, (*impl).parseService)
}

func TestSyncRpc(t *testing.T) {
	tests := []TestCase[ast.Rpc]{
		{
			name: internal.CaseName("sync", true, "rpc"),
			expectedObj: &ast.Rpc{
				ID:         13,
				Name:       ast.Identifier{ID: 1},
				InputType:  ast.Identifier{ID: 3},
				OutputType: ast.Identifier{ID: 7},
			},
			expectedErrs: []error{
				gotUnexpected(&token.Token{ID: 10, Kind: token.KindInt}, token.KindOption),
			},

			content: "rpc MethodName (Request) returns (Response) { 1 }",
			indices: "a--bc---------def------ghi------jkl-------mnopqrst",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'e', 'f'}, {'f', 'g'},
				{'g', 'h'}, {'i', 'j'}, {'k', 'l'}, {'l', 'm'},
				{'m', 'n'}, {'o', 'p'}, {'q', 'r'}, {'s', 't'},
			},
			kinds: []token.Kind{
				token.KindIdentifier,
				token.KindIdentifier, // MethodName
				token.KindLeftParen,
				token.KindIdentifier, // Request
				token.KindRightParen,
				token.KindIdentifier, // returns
				token.KindLeftParen,
				token.KindIdentifier, // Response
				token.KindRightParen,
				token.KindLeftBrace,
				token.KindInt,
				token.KindRightBrace,
			},
		},
	}

	runTestCases(t, tests, checkRpc, (*impl).parseRpc)
}

func TestSyncExtend(t *testing.T) {
	tests := []TestCase[ast.Extend]{
		{
			name: internal.CaseName("sync", true, "extend"),
			expectedObj: &ast.Extend{
				ID:   6,
				Name: ast.Identifier{ID: 1},
			},
			expectedErrs: []error{
				gotUnexpected(&token.Token{ID: 3, Kind: token.KindStr}, token.KindOption, token.KindField),
			},

			content: "extend Test { 'a' }",
			indices: "a-----bc---defg--hij",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'e', 'f'},
				{'g', 'h'}, {'i', 'j'},
			},
			kinds: []token.Kind{
				token.KindIdentifier,
				token.KindIdentifier, // Test
				token.KindLeftBrace,
				token.KindStr,
				token.KindRightBrace,
			},
		},
	}

	runTestCases(t, tests, checkExtend, (*impl).parseExtend)
}

// TODO error in text message list
