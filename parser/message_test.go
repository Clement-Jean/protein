package parser

import (
	"testing"

	"github.com/Clement-Jean/protein/ast"
	"github.com/Clement-Jean/protein/parser/internal"
	"github.com/Clement-Jean/protein/token"
)

func checkField(t *testing.T, got ast.Field, expected ast.Field) {
	if got.Type != expected.Type {
		t.Fatalf("expected type %d, got %d", expected.Type, got.Type)
	}

	checkIDs(t, got.ID, expected.ID)
	checkIDs(t, got.LabelID, expected.LabelID)
	checkIDs(t, got.TypeID, expected.TypeID)
	checkIDs(t, got.Name.ID, expected.Name.ID)
	checkIDs(t, got.OptionsID, expected.OptionsID)

	if len(got.Options) != len(expected.Options) {
		t.Fatalf("expected %d options, got %d", len(expected.Options), len(got.Options))
	}
	for j, opt := range got.Options {
		checkOption(t, opt, expected.Options[j])
	}
}

func checkMessage(t *testing.T, got ast.Message, expected ast.Message) {
	checkIDs(t, got.ID, expected.ID)
	checkIDs(t, got.Name.ID, expected.Name.ID)

	if len(got.Options) != len(expected.Options) {
		t.Fatalf("expected %d options, got %d", len(expected.Options), len(got.Options))
	}
	for i, opt := range got.Options {
		checkOption(t, opt, expected.Options[i])
	}

	if len(got.ReservedTags) != len(expected.ReservedTags) {
		t.Fatalf("expected %d reserved tags, got %d", len(expected.ReservedTags), len(got.ReservedTags))
	}
	for i, opt := range got.ReservedTags {
		checkReservedTags(t, opt, expected.ReservedTags[i])
	}
	if len(got.ReservedNames) != len(expected.ReservedNames) {
		t.Fatalf("expected %d reserved names, got %d", len(expected.ReservedNames), len(got.ReservedNames))
	}
	for i, opt := range got.ReservedNames {
		checkReservedNames(t, opt, expected.ReservedNames[i])
	}

	if len(got.Fields) != len(expected.Fields) {
		t.Fatalf("expected %d field, got %d", len(expected.Fields), len(got.Fields))
	}
	for i, opt := range got.Fields {
		checkField(t, opt, expected.Fields[i])
	}

	if len(got.ExtensionRanges) != len(expected.ExtensionRanges) {
		t.Fatalf("expected %d range, got %d", len(expected.ExtensionRanges), len(got.ExtensionRanges))
	}
	for j, item := range got.ExtensionRanges {
		checkExtensionRange(t, item, expected.ExtensionRanges[j])
	}

	if len(got.Extensions) != len(expected.Extensions) {
		t.Fatalf("expected %d extensions, got %d", len(expected.Extensions), len(got.Extensions))
	}
	for j, item := range got.Extensions {
		checkExtend(t, item, expected.Extensions[j])
	}
}

func TestParseMessage(t *testing.T) {
	tests := []TestCase[ast.Message]{
		{
			name:        internal.CaseName("message", true),
			expectedObj: &ast.Message{ID: 5, Name: ast.Identifier{ID: 1}},

			content: "message Test {}",
			indices: "a------bc---defg",
			locs:    [][2]rune{{'a', 'b'}, {'c', 'd'}, {'e', 'f'}, {'f', 'g'}},
			kinds: []token.Kind{
				token.KindIdentifier,
				token.KindIdentifier, // Test
				token.KindLeftBrace,
				token.KindRightBrace,
			},
		},
		{
			name:        internal.CaseName("message", true, "empty_statement"),
			expectedObj: &ast.Message{ID: 6, Name: ast.Identifier{ID: 1}},

			content: "message Test {;}",
			indices: "a------bc---defgh",
			locs:    [][2]rune{{'a', 'b'}, {'c', 'd'}, {'e', 'f'}, {'f', 'g'}, {'g', 'h'}},
			kinds: []token.Kind{
				token.KindIdentifier,
				token.KindIdentifier, // Test
				token.KindLeftBrace,
				token.KindSemicolon,
				token.KindRightBrace,
			},
		},
		{
			name: internal.CaseName("message", true, "field"),
			expectedObj: &ast.Message{
				ID:   11,
				Name: ast.Identifier{ID: 1},
				Fields: []ast.Field{
					{ID: 10, Type: ast.FieldTypeUint64, TypeID: 3, Name: ast.Identifier{ID: 4}, Tag: ast.Integer{ID: 6}},
				},
			},

			content: "message Test { uint64 id = 1; }",
			indices: "a------bc---defg-----hi-jklmnopq",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'e', 'f'}, {'g', 'h'},
				{'i', 'j'}, {'k', 'l'}, {'m', 'n'}, {'n', 'o'},
				{'p', 'q'},
			},
			kinds: []token.Kind{
				token.KindIdentifier,
				token.KindIdentifier, // Test
				token.KindLeftBrace,
				token.KindIdentifier, // uint64
				token.KindIdentifier, // id
				token.KindEqual,
				token.KindInt,
				token.KindSemicolon,
				token.KindRightBrace,
			},
		},
		{
			name: internal.CaseName("message", true, "map_field"),
			expectedObj: &ast.Message{
				ID:   17,
				Name: ast.Identifier{ID: 1},
				Fields: []ast.Field{
					{
						ID:   16,
						Type: ast.FieldTypeMessage, TypeID: 15,
						Name: ast.Identifier{ID: 9},
						Tag:  ast.Integer{ID: 11},
					},
				},
			},

			content: "message Test { map<string,uint64> id = 1; }",
			indices: "a------bc---defg--hi-----jk-----lmn-opqrstuv",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'e', 'f'}, {'g', 'h'},
				{'h', 'i'}, {'i', 'j'}, {'j', 'k'}, {'k', 'l'},
				{'l', 'm'}, {'n', 'o'}, {'p', 'q'}, {'r', 's'},
				{'s', 't'}, {'u', 'v'},
			},
			kinds: []token.Kind{
				token.KindIdentifier,
				token.KindIdentifier, // Test
				token.KindLeftBrace,
				token.KindIdentifier, // map
				token.KindLeftAngle,
				token.KindIdentifier, // string
				token.KindComma,
				token.KindIdentifier, //uint64
				token.KindRightAngle,
				token.KindIdentifier, // id
				token.KindEqual,
				token.KindInt,
				token.KindSemicolon,
				token.KindRightBrace,
			},
		},
		{
			name: internal.CaseName("message", true, "oneof"),
			expectedObj: &ast.Message{
				ID:   22,
				Name: ast.Identifier{ID: 1},
				Oneofs: []ast.Oneof{{
					ID:   21,
					Name: ast.Identifier{ID: 4},
					Fields: []ast.Field{
						{ID: 19, Type: ast.FieldTypeUint64, TypeID: 7, Name: ast.Identifier{ID: 7}, Tag: ast.Integer{ID: 9}},
						{ID: 20, Type: ast.FieldTypeString, TypeID: 12, Name: ast.Identifier{ID: 12}, Tag: ast.Integer{ID: 14}},
					},
				}},
			},

			content: "message Test { oneof Test2 { uint64 id = 1; string uuid = 2; } }",
			indices: "a------bc---defg----hi----jklm-----no-pqrstuv-----wx---yz12345678",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'e', 'f'}, {'g', 'h'},
				{'i', 'j'}, {'k', 'l'}, {'m', 'n'}, {'o', 'p'},
				{'q', 'r'}, {'s', 't'}, {'t', 'u'}, {'v', 'w'},
				{'x', 'y'}, {'z', '1'}, {'2', '3'}, {'3', '4'},
				{'5', '6'}, {'7', '8'},
			},
			kinds: []token.Kind{
				token.KindIdentifier,
				token.KindIdentifier, // Test
				token.KindLeftBrace,
				token.KindIdentifier, // oneof
				token.KindIdentifier, // Test2
				token.KindLeftBrace,
				token.KindIdentifier, // uint64
				token.KindIdentifier, // id
				token.KindEqual,
				token.KindInt,
				token.KindSemicolon,
				token.KindIdentifier, // string
				token.KindIdentifier, // uuid
				token.KindEqual,
				token.KindInt,
				token.KindSemicolon,
				token.KindRightBrace,
				token.KindRightBrace,
			},
		},
		{
			name: internal.CaseName("message", true, "reserved_tag"),
			expectedObj: &ast.Message{
				ID:   8,
				Name: ast.Identifier{ID: 1},
				ReservedTags: []ast.ReservedTags{
					{ID: 4, Items: []ast.Range{{ID: 4, Start: ast.Integer{ID: 4}, End: ast.Integer{ID: 4}}}},
				},
			},

			content: "message Test { reserved 1; }",
			indices: "a------bc---defg-------hijklm",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'e', 'f'}, {'g', 'h'},
				{'i', 'j'}, {'j', 'k'}, {'l', 'm'},
			},
			kinds: []token.Kind{
				token.KindIdentifier,
				token.KindIdentifier, // Test
				token.KindLeftBrace,
				token.KindIdentifier, // reserved
				token.KindInt,
				token.KindSemicolon,
				token.KindRightBrace,
			},
		},
		{
			name: internal.CaseName("message", true, "reserved_name"),
			expectedObj: &ast.Message{
				ID:   8,
				Name: ast.Identifier{ID: 1},
				ReservedNames: []ast.ReservedNames{
					{ID: 4, Items: []ast.String{{ID: 4}}},
				},
			},

			content: "message Test { reserved '1'; }",
			indices: "a------bc---defg-------hi--jklm",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'e', 'f'}, {'g', 'h'},
				{'i', 'j'}, {'j', 'k'}, {'l', 'm'},
			},
			kinds: []token.Kind{
				token.KindIdentifier,
				token.KindIdentifier, // Test
				token.KindLeftBrace,
				token.KindIdentifier, // reserved
				token.KindStr,
				token.KindSemicolon,
				token.KindRightBrace,
			},
		},
		{
			name: internal.CaseName("message", true, "option"),
			expectedObj: &ast.Message{
				ID: 11, Name: ast.Identifier{ID: 1},
				Options: []ast.Option{{
					ID: 10, Name: ast.Identifier{ID: 4}, Value: ast.Boolean{ID: 6},
				}},
			},

			content: "message Test { option deprecated = true; }",
			indices: "a------bc---defg-----hi---------jklm---nopq",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'e', 'f'}, {'g', 'h'},
				{'i', 'j'}, {'k', 'l'}, {'m', 'n'}, {'n', 'o'},
				{'p', 'q'},
			},
			kinds: []token.Kind{
				token.KindIdentifier,
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
			name: internal.CaseName("message", true, "nested_message"),
			expectedObj: &ast.Message{
				ID:   10,
				Name: ast.Identifier{ID: 1},
				Messages: []ast.Message{
					{ID: 9, Name: ast.Identifier{ID: 4}},
				},
			},

			content: "message Test { message Test2 { } }",
			indices: "a------bc---defg------hi----jklmnop",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'e', 'f'}, {'g', 'h'},
				{'i', 'j'}, {'k', 'l'}, {'m', 'n'}, {'o', 'p'},
			},
			kinds: []token.Kind{
				token.KindIdentifier,
				token.KindIdentifier, // Test
				token.KindLeftBrace,
				token.KindIdentifier, // message
				token.KindIdentifier, // Test2
				token.KindLeftBrace,
				token.KindRightBrace,
				token.KindRightBrace,
			},
		},
		{
			name: internal.CaseName("message", true, "nested_enum"),
			expectedObj: &ast.Message{
				ID:   10,
				Name: ast.Identifier{ID: 1},
				Messages: []ast.Message{
					{ID: 9, Name: ast.Identifier{ID: 4}},
				},
			},

			content: "message Test { enum Test2 { } }",
			indices: "a------bc---defg---hi----jklmnop",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'e', 'f'}, {'g', 'h'},
				{'i', 'j'}, {'k', 'l'}, {'m', 'n'}, {'o', 'p'},
			},
			kinds: []token.Kind{
				token.KindIdentifier, // message
				token.KindIdentifier, // Test
				token.KindLeftBrace,
				token.KindIdentifier, // enum
				token.KindIdentifier, // Test2
				token.KindLeftBrace,
				token.KindRightBrace,
				token.KindRightBrace,
			},
		},
		{
			name: internal.CaseName("message", true, "nested_extend"),
			expectedObj: &ast.Message{
				ID:   10,
				Name: ast.Identifier{ID: 1},
				Extensions: []ast.Extend{
					{ID: 9, Name: ast.Identifier{ID: 4}},
				},
			},

			content: "message Test { extend Test2 {} }",
			indices: "a------bc---defg-----hi----jklmno",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'e', 'f'}, {'g', 'h'},
				{'i', 'j'}, {'k', 'l'}, {'l', 'm'}, {'n', 'o'},
			},
			kinds: []token.Kind{
				token.KindIdentifier, // message
				token.KindIdentifier, // Test
				token.KindLeftBrace,
				token.KindIdentifier, // extend
				token.KindIdentifier, // Test2
				token.KindLeftBrace,
				token.KindRightBrace,
				token.KindRightBrace,
			},
		},
		{
			name: internal.CaseName("message", true, "extension_range"),
			expectedObj: &ast.Message{
				ID:   12,
				Name: ast.Identifier{ID: 1},
				ExtensionRanges: []ast.ExtensionRange{
					{
						ID: 11,
						Ranges: []ast.Range{
							{ID: 10, Start: ast.Integer{ID: 4}, End: ast.Integer{ID: 6}},
						},
					},
				},
			},

			content: "message Test { extensions 1000 to max; }",
			indices: "a------bc---defg---------hi---jk-lm--nopq",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'e', 'f'}, {'g', 'h'},
				{'i', 'j'}, {'k', 'l'}, {'m', 'n'}, {'n', 'o'},
				{'p', 'q'},
			},
			kinds: []token.Kind{
				token.KindIdentifier,
				token.KindIdentifier, // Test
				token.KindLeftBrace,
				token.KindIdentifier, // extensions
				token.KindInt,
				token.KindIdentifier, // to
				token.KindIdentifier, // max
				token.KindSemicolon,
				token.KindRightBrace,
			},
		},
		{
			name: internal.CaseName("message", false, "unexpected_int"),
			expectedObj: &ast.Message{
				ID:   6,
				Name: ast.Identifier{ID: 1},
			},
			expectedErrs: []error{
				gotUnexpected(&token.Token{ID: 3, Kind: token.KindInt}, token.KindOption, token.KindReserved, token.KindField),
			},

			content: "message Test { 2 }",
			indices: "a------bc---defghij",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'e', 'f'}, {'g', 'h'},
				{'i', 'j'},
			},
			kinds: []token.Kind{
				token.KindIdentifier,
				token.KindIdentifier, // Test
				token.KindLeftBrace,
				token.KindInt,
				token.KindRightBrace,
			},
		},
		{
			name: internal.CaseName("message", false, "expected_identifier"),
			expectedErrs: []error{
				gotUnexpected(&token.Token{ID: 1, Kind: token.KindLeftBrace}, token.KindIdentifier),
			},

			content: "message {}",
			indices: "a------bcde",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'d', 'e'},
			},
			kinds: []token.Kind{
				token.KindIdentifier, // message
				token.KindLeftBrace,
				token.KindRightBrace,
			},
		},
		{
			name: internal.CaseName("message", false, "expected_left_brace"),
			expectedErrs: []error{
				gotUnexpected(&token.Token{ID: 2, Kind: token.KindLeftSquare}, token.KindLeftBrace),
			},

			content: "message Test [}",
			indices: "a------bc---defg",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'d', 'e'}, {'f', 'g'},
			},
			kinds: []token.Kind{
				token.KindIdentifier,
				token.KindIdentifier,
				token.KindLeftSquare,
				token.KindRightBrace,
			},
		},
		{
			name: internal.CaseName("message", false, "expected_right_brace"),
			expectedErrs: []error{
				gotUnexpected(&token.Token{ID: 3, Kind: token.KindEOF}, token.KindRightBrace),
			},

			content: "message Test {",
			indices: "a------bc---def",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'e', 'f'},
			},
			kinds: []token.Kind{
				token.KindIdentifier,
				token.KindIdentifier, // Test
				token.KindLeftBrace,
			},
		},
	}

	wrap := func(p *impl) (ast.Message, []error) { return p.parseMessage(1) }
	runTestCases(t, tests, checkMessage, wrap)
}

func TestParseField(t *testing.T) {
	tests := []TestCase[ast.Field]{
		{
			keepFirstToken: true,
			name:           internal.CaseName("field", true, "type_uint64"),
			expectedObj: &ast.Field{
				ID: 6, Type: ast.FieldTypeUint64, TypeID: 0, Name: ast.Identifier{ID: 1}, Tag: ast.Integer{ID: 3},
			},

			content: "uint64 id = 1;",
			indices: "a-----bc-defghi",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'e', 'f'}, {'g', 'h'},
				{'h', 'i'},
			},
			kinds: []token.Kind{
				token.KindIdentifier, // uint64
				token.KindIdentifier, // id
				token.KindEqual,
				token.KindInt,
				token.KindSemicolon,
			},
		},
		{
			keepFirstToken: true,
			name:           internal.CaseName("field", true, "type_unknown"),
			expectedObj: &ast.Field{
				ID: 6, Type: ast.FieldTypeUnknown, TypeID: 0, Name: ast.Identifier{ID: 1}, Tag: ast.Integer{ID: 3},
			},

			content: "Test id = 1;",
			indices: "a---bc-defghi",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'e', 'f'}, {'g', 'h'},
				{'h', 'i'},
			},
			kinds: []token.Kind{
				token.KindIdentifier, // Test
				token.KindIdentifier, // id
				token.KindEqual,
				token.KindInt,
				token.KindSemicolon,
			},
		},
		{
			keepFirstToken: true,
			name:           internal.CaseName("field", true, "label"),
			expectedObj: &ast.Field{
				ID:    7,
				Label: ast.FieldLabelRepeated, LabelID: 0,
				Type: ast.FieldTypeUint64, TypeID: 1,
				Name: ast.Identifier{ID: 2},
				Tag:  ast.Integer{ID: 4},
			},

			content: "repeated uint64 ids = 1;",
			indices: "a-------bc-----de--fghijk",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'e', 'f'}, {'g', 'h'},
				{'i', 'j'}, {'j', 'k'},
			},
			kinds: []token.Kind{
				token.KindIdentifier, // repeated
				token.KindIdentifier, // uint64
				token.KindIdentifier, // ids
				token.KindEqual,
				token.KindInt,
				token.KindSemicolon,
			},
		},
		{
			keepFirstToken: true,
			name:           internal.CaseName("field", true, "option"),
			expectedObj: &ast.Field{
				ID:   13,
				Type: ast.FieldTypeUint64, TypeID: 0,
				Name: ast.Identifier{ID: 1},
				Tag:  ast.Integer{ID: 3},
				Options: []ast.Option{
					{ID: 11, Name: ast.Identifier{ID: 5}, Value: ast.Boolean{ID: 7}},
				}, OptionsID: 12,
			},

			content: "uint64 id = 1 [deprecated = true];",
			indices: "a-----bc-defghij---------klmn---opq",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'e', 'f'}, {'g', 'h'},
				{'i', 'j'}, {'j', 'k'}, {'l', 'm'}, {'n', 'o'},
				{'o', 'p'}, {'p', 'q'},
			},
			kinds: []token.Kind{
				token.KindIdentifier, // uint64
				token.KindIdentifier, // id
				token.KindEqual,
				token.KindInt,
				token.KindLeftSquare,
				token.KindIdentifier, // deprecated
				token.KindEqual,
				token.KindIdentifier, // true
				token.KindRightSquare,
				token.KindSemicolon,
			},
		},
		{
			keepFirstToken: true,
			name:           internal.CaseName("field", true, "empty_option"),
			expectedObj: &ast.Field{
				ID:   9,
				Type: ast.FieldTypeUint64, TypeID: 0,
				Name:      ast.Identifier{ID: 1},
				Tag:       ast.Integer{ID: 3},
				OptionsID: 8,
			},

			content: "uint64 id = 1 [];",
			indices: "a-----bc-defghijkl",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'e', 'f'}, {'g', 'h'},
				{'i', 'j'}, {'j', 'k'}, {'k', 'l'},
			},
			kinds: []token.Kind{
				token.KindIdentifier, // uint64
				token.KindIdentifier, // id
				token.KindEqual,
				token.KindInt,
				token.KindLeftSquare,
				token.KindRightSquare,
				token.KindSemicolon,
			},
		},
		{
			keepFirstToken: true,
			name:           internal.CaseName("field", true, "options"),
			expectedObj: &ast.Field{
				ID:        30,
				TypeID:    0,
				Type:      ast.FieldTypeUint64,
				Name:      ast.Identifier{ID: 1},
				Tag:       ast.Integer{ID: 3},
				OptionsID: 29,
				Options: []ast.Option{
					{
						ID: 25, Name: ast.Identifier{ID: 5},
						Value: ast.TextMessage{
							ID: 24,
							Fields: []ast.TextField{
								{
									ID:    23,
									Name:  ast.Identifier{ID: 8},
									Value: ast.Identifier{ID: 10},
								},
							},
						},
					},
					{
						ID: 28, Name: ast.Identifier{ID: 13},
						Value: ast.TextMessage{
							ID: 27,
							Fields: []ast.TextField{
								{
									ID:    26,
									Name:  ast.Identifier{ID: 16},
									Value: ast.Identifier{ID: 18},
								},
							},
						},
					},
				},
			},

			content: "uint64 id = 1 [a = { b: c }, d = { e: f }];",
			indices: "a-----bc-defghijklmnopqrstuvwxyz123456789ABC",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'e', 'f'}, {'g', 'h'},
				{'i', 'j'}, {'j', 'k'}, {'l', 'm'}, {'n', 'o'},
				{'p', 'q'}, {'q', 'r'}, {'s', 't'}, {'u', 'v'},
				{'v', 'w'}, {'x', 'y'}, {'z', '1'}, {'2', '3'},
				{'4', '5'}, {'5', '6'}, {'7', '8'}, {'9', 'A'},
				{'A', 'B'}, {'B', 'C'},
			},
			kinds: []token.Kind{
				token.KindIdentifier, // uint64
				token.KindIdentifier, // id
				token.KindEqual,
				token.KindInt,
				token.KindLeftSquare,
				token.KindIdentifier, // a
				token.KindEqual,
				token.KindLeftBrace,
				token.KindIdentifier, // b
				token.KindColon,
				token.KindIdentifier, // c
				token.KindRightBrace,
				token.KindComma,
				token.KindIdentifier, // d
				token.KindEqual,
				token.KindLeftBrace,
				token.KindIdentifier, // e
				token.KindColon,
				token.KindIdentifier, // f
				token.KindRightBrace,
				token.KindRightSquare,
				token.KindSemicolon,
			},
		},
		{
			keepFirstToken: true,
			name:           internal.CaseName("field", false, "expected_identifier_after_label"),
			expectedErrs: []error{
				gotUnexpected(&token.Token{ID: 1, Kind: token.KindStr}, token.KindIdentifier),
			},

			content: "repeated 'uint64' ids = 1;",
			indices: "a-------bc-------de--fghijk",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'e', 'f'}, {'g', 'h'},
				{'i', 'j'}, {'j', 'k'},
			},
			kinds: []token.Kind{
				token.KindIdentifier, // repeated
				token.KindStr,
				token.KindIdentifier, // ids
				token.KindEqual,
				token.KindInt,
				token.KindSemicolon,
			},
		},
		{
			keepFirstToken: true,
			name:           internal.CaseName("field", false, "expected_identifier_after_type"),
			expectedErrs: []error{
				gotUnexpected(&token.Token{ID: 1, Kind: token.KindStr}, token.KindIdentifier),
			},

			content: "uint64 'id' = 1;",
			indices: "a-----bc---defghi",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'e', 'f'}, {'g', 'h'},
				{'h', 'i'},
			},
			kinds: []token.Kind{
				token.KindIdentifier, // uint64
				token.KindStr,
				token.KindEqual,
				token.KindInt,
				token.KindSemicolon,
			},
		},
		{
			keepFirstToken: true,
			name:           internal.CaseName("field", false, "expected_equal_after_name"),
			expectedErrs: []error{
				gotUnexpected(&token.Token{ID: 2, Kind: token.KindInt}, token.KindEqual),
			},

			content: "uint64 id 1;",
			indices: "a-----bc-defg",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'e', 'f'}, {'f', 'g'},
			},
			kinds: []token.Kind{
				token.KindIdentifier, // uint64
				token.KindIdentifier, // id
				token.KindInt,
				token.KindSemicolon,
			},
		},
		{
			keepFirstToken: true,
			name:           internal.CaseName("field", false, "expected_int_after_equal"),
			expectedErrs: []error{
				gotUnexpected(&token.Token{ID: 3, Kind: token.KindStr}, token.KindInt),
			},

			content: "uint64 id = '1';",
			indices: "a-----bc-defg--hi",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'e', 'f'}, {'g', 'h'},
				{'h', 'i'},
			},
			kinds: []token.Kind{
				token.KindIdentifier, // uint64
				token.KindIdentifier, // id
				token.KindEqual,
				token.KindStr,
				token.KindSemicolon,
			},
		},
		{
			keepFirstToken: true,
			name:           internal.CaseName("field", false, "expected_semicolon"),
			expectedErrs: []error{
				gotUnexpected(&token.Token{ID: 4, Kind: token.KindEOF}, token.KindSemicolon),
			},

			content: "uint64 id = 1",
			indices: "a-----bc-defgh",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'e', 'f'}, {'g', 'h'},
			},
			kinds: []token.Kind{
				token.KindIdentifier, // uint64
				token.KindIdentifier, // id
				token.KindEqual,
				token.KindInt,
			},
		},
		{
			keepFirstToken: true,
			name:           internal.CaseName("field", false, "expected_right_square"),
			expectedErrs: []error{
				gotUnexpected(&token.Token{ID: 8, Kind: token.KindEOF}, token.KindRightSquare),
				gotUnexpected(&token.Token{ID: 8, Kind: token.KindEOF}, token.KindSemicolon),
			},

			content: "uint64 id = 1 [deprecated = true",
			indices: "a-----bc-defghij---------klmn---o",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'e', 'f'}, {'g', 'h'},
				{'i', 'j'}, {'j', 'k'}, {'l', 'm'}, {'n', 'o'},
			},
			kinds: []token.Kind{
				token.KindIdentifier, // uint64
				token.KindIdentifier, // id
				token.KindEqual,
				token.KindInt,
				token.KindLeftSquare,
				token.KindIdentifier,
				token.KindEqual,
				token.KindIdentifier, // true
			},
		},
	}

	runTestCases(t, tests, checkField, (*impl).parseField)
}

func TestParseMapField(t *testing.T) {
	tests := []TestCase[ast.Field]{
		{
			keepFirstToken: true,
			name:           internal.CaseName("map_field", true),
			expectedObj: &ast.Field{
				ID: 12, Type: ast.FieldTypeMessage, TypeID: 11, Name: ast.Identifier{ID: 6}, Tag: ast.Integer{ID: 8},
			},

			content: "map<string, uint64> ids = 1;",
			indices: "a--bc-----def-----ghi--jklmno",
			locs: [][2]rune{
				{'a', 'b'}, {'b', 'c'}, {'c', 'd'}, {'d', 'e'},
				{'f', 'g'}, {'g', 'h'}, {'i', 'j'}, {'k', 'l'},
				{'m', 'n'}, {'n', 'o'},
			},
			kinds: []token.Kind{
				token.KindIdentifier, // map
				token.KindLeftAngle,
				token.KindIdentifier, // string
				token.KindComma,
				token.KindIdentifier, // uint64
				token.KindRightAngle,
				token.KindIdentifier, // ids
				token.KindEqual,
				token.KindInt,
				token.KindSemicolon,
			},
		},
		{
			keepFirstToken: true,
			name:           internal.CaseName("map_field", false, "expected_left_angle"),
			expectedErrs: []error{
				gotUnexpected(&token.Token{ID: 1, Kind: token.KindLeftSquare}, token.KindLeftAngle),
			},

			content: "map[string, uint64> ids = 1;",
			indices: "a--bc-----def-----ghi--jklmno",
			locs: [][2]rune{
				{'a', 'b'}, {'b', 'c'}, {'c', 'd'}, {'d', 'e'},
				{'f', 'g'}, {'g', 'h'}, {'i', 'j'}, {'k', 'l'},
				{'m', 'n'}, {'n', 'o'},
			},
			kinds: []token.Kind{
				token.KindIdentifier, // map
				token.KindLeftSquare,
				token.KindIdentifier, // string
				token.KindComma,
				token.KindIdentifier, // uint64
				token.KindRightAngle,
				token.KindIdentifier, // ids
				token.KindEqual,
				token.KindInt,
				token.KindSemicolon,
			},
		},
		{
			keepFirstToken: true,
			name:           internal.CaseName("map_field", false, "expected_key_identifier"),
			expectedErrs: []error{
				gotUnexpected(&token.Token{ID: 2, Kind: token.KindStr}, token.KindIdentifier),
			},

			content: "map<'string', uint64> ids = 1;",
			indices: "a--bc-------def-----ghi--jklmno",
			locs: [][2]rune{
				{'a', 'b'}, {'b', 'c'}, {'c', 'd'}, {'d', 'e'},
				{'f', 'g'}, {'g', 'h'}, {'i', 'j'}, {'k', 'l'},
				{'m', 'n'}, {'n', 'o'},
			},
			kinds: []token.Kind{
				token.KindIdentifier, // map
				token.KindLeftAngle,
				token.KindStr,
				token.KindComma,
				token.KindIdentifier, // uint64
				token.KindRightAngle,
				token.KindIdentifier, // ids
				token.KindEqual,
				token.KindInt,
				token.KindSemicolon,
			},
		},
		{
			keepFirstToken: true,
			name:           internal.CaseName("map_field", false, "expected_comma"),
			expectedErrs: []error{
				gotUnexpected(&token.Token{ID: 3, Kind: token.KindIdentifier}, token.KindComma),
			},

			content: "map<string uint64> ids = 1;",
			indices: "a--bc-----de-----fgh--ijklmn",
			locs: [][2]rune{
				{'a', 'b'}, {'b', 'c'}, {'c', 'd'}, {'e', 'f'},
				{'f', 'g'}, {'h', 'i'}, {'j', 'k'}, {'l', 'm'},
				{'m', 'n'},
			},
			kinds: []token.Kind{
				token.KindIdentifier, // map
				token.KindLeftAngle,
				token.KindIdentifier, // string
				token.KindIdentifier, // uint64
				token.KindRightAngle,
				token.KindIdentifier, // ids
				token.KindEqual,
				token.KindInt,
				token.KindSemicolon,
			},
		},
		{
			keepFirstToken: true,
			name:           internal.CaseName("map_field", false, "expected_Value_identifier"),
			expectedErrs: []error{
				gotUnexpected(&token.Token{ID: 4, Kind: token.KindStr}, token.KindIdentifier),
			},

			content: "map<string, 'uint64'> ids = 1;",
			indices: "a--bc-----def-------ghi--jklmno",
			locs: [][2]rune{
				{'a', 'b'}, {'b', 'c'}, {'c', 'd'}, {'d', 'e'},
				{'f', 'g'}, {'g', 'h'}, {'i', 'j'}, {'k', 'l'},
				{'m', 'n'}, {'n', 'o'},
			},
			kinds: []token.Kind{
				token.KindIdentifier, // map
				token.KindLeftAngle,
				token.KindIdentifier, // string
				token.KindComma,
				token.KindStr,
				token.KindRightAngle,
				token.KindIdentifier, // ids
				token.KindEqual,
				token.KindInt,
				token.KindSemicolon,
			},
		},
		{
			keepFirstToken: true,
			name:           internal.CaseName("map_field", false, "expected_right_angle"),
			expectedErrs: []error{
				gotUnexpected(&token.Token{ID: 5, Kind: token.KindRightSquare}, token.KindRightAngle),
			},

			content: "map<string, uint64] ids = 1;",
			indices: "a--bc-----def-----ghi--jklmno",
			locs: [][2]rune{
				{'a', 'b'}, {'b', 'c'}, {'c', 'd'}, {'d', 'e'},
				{'f', 'g'}, {'g', 'h'}, {'i', 'j'}, {'k', 'l'},
				{'m', 'n'}, {'n', 'o'},
			},
			kinds: []token.Kind{
				token.KindIdentifier, // map
				token.KindLeftAngle,
				token.KindIdentifier, // string
				token.KindComma,
				token.KindIdentifier, // uint64
				token.KindRightSquare,
				token.KindIdentifier, // ids
				token.KindEqual,
				token.KindInt,
				token.KindSemicolon,
			},
		},
		{
			keepFirstToken: true,
			name:           internal.CaseName("map_field", false, "expected_right_square"),
			expectedErrs: []error{
				gotUnexpected(&token.Token{ID: 13, Kind: token.KindEOF}, token.KindRightSquare),
				gotUnexpected(&token.Token{ID: 13, Kind: token.KindEOF}, token.KindSemicolon),
			},

			content: "map<string, uint64> ids = 1 [deprecated = true",
			indices: "a--bc-----def-----ghi--jklmnop---------qrst---u",
			locs: [][2]rune{
				{'a', 'b'}, {'b', 'c'}, {'c', 'd'}, {'d', 'e'},
				{'f', 'g'}, {'g', 'h'}, {'i', 'j'}, {'k', 'l'},
				{'m', 'n'}, {'o', 'p'}, {'p', 'q'}, {'r', 's'},
				{'t', 'u'},
			},
			kinds: []token.Kind{
				token.KindIdentifier, // map
				token.KindLeftAngle,
				token.KindIdentifier, // string
				token.KindComma,
				token.KindIdentifier, // uint64
				token.KindRightAngle,
				token.KindIdentifier, // ids
				token.KindEqual,
				token.KindInt,
				token.KindLeftSquare,
				token.KindIdentifier, // deprecated
				token.KindEqual,
				token.KindIdentifier, // true
			},
		},
	}

	runTestCases(t, tests, checkField, (*impl).parseMapField)
}
