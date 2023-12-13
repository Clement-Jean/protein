package parser

import (
	"testing"

	"github.com/Clement-Jean/protein/ast"
	"github.com/Clement-Jean/protein/parser/internal"
	"github.com/Clement-Jean/protein/token"
)

func checkTextField(t *testing.T, got, expected ast.TextField) {
	checkIDs(t, got.ID, expected.ID)
	checkIdentifier(t, got.Name, expected.Name)

	switch any(expected.Value).(type) {
	case ast.TextScalarList:
		gotValueType, ok := got.Value.(ast.TextScalarList)

		if !ok {
			t.Fatalf("expected scalar list, got %s", got.Value)
		}

		expectedValueType := expected.Value.(ast.TextScalarList)
		checkIDs(t, gotValueType.ID, expectedValueType.ID)

		if len(gotValueType.Values) != len(expectedValueType.Values) {
			t.Fatalf("expected %d values, got %d", len(expectedValueType.Values), len(gotValueType.Values))
		}

		for i, value := range expectedValueType.Values {
			checkIDs(t, gotValueType.Values[i].GetID(), value.GetID())
		}
	case ast.TextMessageList:
		gotValueType, ok := got.Value.(ast.TextMessageList)

		if !ok {
			t.Fatalf("expected scalar list, got %s", got.Value)
		}

		expectedValueType := expected.Value.(ast.TextMessageList)
		checkIDs(t, gotValueType.ID, expectedValueType.ID)

		if len(gotValueType.Values) != len(expectedValueType.Values) {
			t.Fatalf("expected %d values, got %d", len(expectedValueType.Values), len(gotValueType.Values))
		}

		for i, value := range expectedValueType.Values {
			checkIDs(t, gotValueType.Values[i].GetID(), value.GetID())
		}
	case ast.Identifier:
		gotValueType, ok := got.Value.(ast.Identifier)

		if !ok {
			t.Fatalf("expected scalar identifier, got %s", got.Value)
		}
		checkIDs(t, gotValueType.ID, expected.Value.(ast.Identifier).ID)
	case ast.Integer:
		gotValueType, ok := got.Value.(ast.Integer)

		if !ok {
			t.Fatalf("expected scalar integer, got %s", got.Value)
		}
		checkIDs(t, gotValueType.ID, expected.Value.(ast.Integer).ID)
	case ast.Boolean:
		gotValueType, ok := got.Value.(ast.Boolean)

		if !ok {
			t.Fatalf("expected scalar boolean, got %s", got.Value)
		}
		checkIDs(t, gotValueType.ID, expected.Value.(ast.Boolean).ID)
	case ast.String:
		gotValueType, ok := got.Value.(ast.String)

		if !ok {
			t.Fatalf("expected scalar string, got %s", got.Value)
		}
		checkIDs(t, gotValueType.ID, expected.Value.(ast.String).ID)
	case ast.Float:
		gotValueType, ok := got.Value.(ast.Float)

		if !ok {
			t.Fatalf("expected scalar float, got %s", got.Value)
		}
		checkIDs(t, gotValueType.ID, expected.Value.(ast.Float).ID)
	case ast.TextMessage:
		gotValueType, ok := got.Value.(ast.TextMessage)

		if !ok {
			t.Fatalf("expected text message, got %s", got.Value)
		}

		expectedValueType := expected.Value.(ast.TextMessage)
		checkIDs(t, gotValueType.ID, expectedValueType.ID)
		if len(gotValueType.Fields) != len(expectedValueType.Fields) {
			t.Fatalf("expected %d fields, got %d", len(expectedValueType.Fields), len(gotValueType.Fields))
		}

		for i, field := range expectedValueType.Fields {
			checkIDs(t, gotValueType.Fields[i].ID, field.ID)
		}
	default:
		t.Fatal("this should never happen!")
	}
}

func checkTextMessage(t *testing.T, got, expected ast.TextMessage) {
	checkIDs(t, got.ID, expected.ID)

	if len(got.Fields) != len(expected.Fields) {
		t.Fatalf("expected %d field, got %d", len(expected.Fields), len(got.Fields))
	}

	for i, field := range got.Fields {
		checkTextField(t, field, expected.Fields[i])
	}
}

func TestParseFullyQualifiedIdentifier(t *testing.T) {
	tests := []TestCase[ast.Identifier]{
		{
			keepFirstToken: true,
			name:           internal.CaseName("fully_qualified_identifier", true),
			expectedErrs: []error{
				gotUnexpected(&token.Token{ID: 2, Kind: token.KindStr}, token.KindIdentifier),
			},

			content: "test.'test'.test",
			indices: "a---bc-----de---f",
			locs:    [][2]rune{{'a', 'b'}, {'b', 'c'}, {'c', 'd'}, {'d', 'e'}, {'e', 'f'}},
			kinds: []token.Kind{
				token.KindIdentifier, // test
				token.KindDot,
				token.KindStr,
				token.KindDot,
				token.KindIdentifier, // test
			},
		},
	}

	runTestCases(t, tests, checkIdentifier, (*impl).parseFullyQualifiedIdentifier)
}

func TestParseTextFieldName(t *testing.T) {
	tests := []TestCase[ast.Identifier]{
		{
			keepFirstToken: true,
			name:           internal.CaseName("text_field_name", true),
			expectedObj:    ast.Identifier{ID: 0},

			content: "reg_scalar",
			indices: "a---------b",
			locs:    [][2]rune{{'a', 'b'}},
			kinds: []token.Kind{
				token.KindIdentifier, //reg_scalar
			},
		},
		{
			keepFirstToken: true,
			name:           internal.CaseName("text_field_name", true, "between_square_brackets"),
			expectedObj:    ast.Identifier{ID: 1},

			content: "[reg_scalar]",
			indices: "ab---------cd",
			locs:    [][2]rune{{'a', 'b'}, {'b', 'c'}, {'c', 'd'}},
			kinds: []token.Kind{
				token.KindLeftSquare,
				token.KindIdentifier, //reg_scalar
				token.KindRightSquare,
			},
		},
		{
			keepFirstToken: true,
			name:           internal.CaseName("text_field_name", true, "fully_qualified_between_square_brackets"),
			expectedObj:    ast.Identifier{ID: 6},

			content: "[reg_scalar.test]",
			indices: "ab---------cd---ef",
			locs:    [][2]rune{{'a', 'b'}, {'b', 'c'}, {'c', 'd'}, {'d', 'e'}, {'e', 'f'}},
			kinds: []token.Kind{
				token.KindLeftSquare,
				token.KindIdentifier, //reg_scalar
				token.KindDot,
				token.KindIdentifier, // test
				token.KindRightSquare,
			},
		},
		{
			keepFirstToken: true,
			name:           internal.CaseName("text_field_name", true, "domain_typename"),
			expectedObj:    ast.Identifier{ID: 16, Parts: []token.UniqueID{14, 15}},

			content: "[type.googleapis.com/com.foo.any]",
			indices: "ab---cd---------ef--gh--ij--kl--mn",
			locs: [][2]rune{
				{'a', 'b'}, {'b', 'c'}, {'c', 'd'}, {'d', 'e'},
				{'e', 'f'}, {'f', 'g'}, {'g', 'h'}, {'h', 'i'},
				{'i', 'j'}, {'j', 'k'}, {'k', 'l'}, {'l', 'm'},
				{'m', 'n'},
			},
			kinds: []token.Kind{
				token.KindLeftSquare,
				token.KindIdentifier, // type
				token.KindDot,
				token.KindIdentifier, // googleapis
				token.KindDot,
				token.KindIdentifier, // com
				token.KindSlash,
				token.KindIdentifier, // com
				token.KindDot,
				token.KindIdentifier, // foo
				token.KindDot,
				token.KindIdentifier, // any
				token.KindRightSquare,
			},
		},
		{
			keepFirstToken: true,
			name:           internal.CaseName("text_field_name", false, "expected_identifier"),
			expectedErrs: []error{
				gotUnexpected(&token.Token{ID: 0, Kind: token.KindInt}, token.KindIdentifier, token.KindLeftSquare),
			},

			content: "2",
			indices: "ab",
			locs:    [][2]rune{{'a', 'b'}},
			kinds: []token.Kind{
				token.KindInt,
			},
		},
		{
			keepFirstToken: true,
			name:           internal.CaseName("text_field_name", false, "expected_identifier_in_square"),
			expectedErrs: []error{
				gotUnexpected(&token.Token{ID: 1, Kind: token.KindStr}, token.KindIdentifier),
			},

			content: "['reg_scalar']",
			indices: "ab-----------cd",
			locs:    [][2]rune{{'a', 'b'}, {'b', 'c'}, {'c', 'd'}},
			kinds: []token.Kind{
				token.KindLeftSquare,
				token.KindStr,
				token.KindRightSquare,
			},
		},
		{
			keepFirstToken: true,
			name:           internal.CaseName("text_field_name", false, "expected_identifier_after_slash"),
			expectedErrs: []error{
				gotUnexpected(&token.Token{ID: 3, Kind: token.KindStr}, token.KindIdentifier),
			},

			content: "[test/'reg_scalar']",
			indices: "ab---cd-----------ef",
			locs:    [][2]rune{{'a', 'b'}, {'b', 'c'}, {'c', 'd'}, {'d', 'e'}, {'e', 'f'}},
			kinds: []token.Kind{
				token.KindLeftSquare,
				token.KindIdentifier,
				token.KindSlash,
				token.KindStr,
				token.KindRightSquare,
			},
		},
		{
			keepFirstToken: true,
			name:           internal.CaseName("text_field_name", false, "expected_right_square"),
			expectedErrs: []error{
				gotUnexpected(&token.Token{ID: 2, Kind: token.KindEOF}, token.KindRightSquare),
			},

			content: "[reg_scalar",
			indices: "ab---------c",
			locs:    [][2]rune{{'a', 'b'}, {'b', 'c'}},
			kinds: []token.Kind{
				token.KindLeftSquare,
				token.KindIdentifier,
			},
		},
	}

	runTestCases(t, tests, checkIdentifier, (*impl).parseTextFieldName)
}

func TestParseTextField(t *testing.T) {
	tests := []TestCase[ast.TextField]{
		{
			keepFirstToken: true,
			name:           internal.CaseName("text_field", true),
			expectedObj:    ast.TextField{ID: 4, Value: ast.Integer{ID: 2}},

			content: "reg_scalar: 10",
			indices: "a---------bcd-e",
			locs:    [][2]rune{{'a', 'b'}, {'b', 'c'}, {'d', 'e'}},
			kinds: []token.Kind{
				token.KindIdentifier, //reg_scalar
				token.KindColon,
				token.KindInt,
			},
		},
		{
			keepFirstToken: true,
			name:           internal.CaseName("text_field", true, "skip_semicolon"),
			expectedObj:    ast.TextField{ID: 5, Value: ast.Integer{ID: 2}},

			content: "reg_scalar: 10;",
			indices: "a---------bcd-ef",
			locs:    [][2]rune{{'a', 'b'}, {'b', 'c'}, {'d', 'e'}, {'e', 'f'}},
			kinds: []token.Kind{
				token.KindIdentifier, //reg_scalar
				token.KindColon,
				token.KindInt,
				token.KindSemicolon,
			},
		},
		{
			keepFirstToken: true,
			name:           internal.CaseName("text_field", true, "list"),
			expectedObj: ast.TextField{ID: 9, Value: ast.TextScalarList{ID: 8, Values: []ast.Expression{
				ast.Integer{ID: 3}, ast.Integer{ID: 5},
			}}},

			content: "reg_scalar: [10, 11]",
			indices: "a---------bcde-fgh-ij",
			locs: [][2]rune{
				{'a', 'b'}, {'b', 'c'}, {'d', 'e'}, {'e', 'f'},
				{'f', 'g'}, {'h', 'i'}, {'i', 'j'},
			},
			kinds: []token.Kind{
				token.KindIdentifier, //reg_scalar
				token.KindColon,
				token.KindLeftSquare,
				token.KindInt,
				token.KindComma,
				token.KindInt,
				token.KindRightSquare,
			},
		},
		{
			keepFirstToken: true,
			name:           internal.CaseName("text_field", true, "message"),
			expectedObj:    ast.TextField{ID: 6, Value: ast.TextMessage{ID: 5}},

			content: "reg_scalar: {}",
			indices: "a---------bcdef",
			locs:    [][2]rune{{'a', 'b'}, {'b', 'c'}, {'d', 'e'}, {'e', 'f'}},
			kinds: []token.Kind{
				token.KindIdentifier, //reg_scalar
				token.KindColon,
				token.KindLeftBrace,
				token.KindRightBrace,
			},
		},
		{
			keepFirstToken: true,
			name:           internal.CaseName("text_field", true, "message_without_colon"),
			expectedObj:    ast.TextField{ID: 5, Value: ast.TextMessage{ID: 4}},

			content: "reg_scalar {}",
			indices: "a---------bcde",
			locs:    [][2]rune{{'a', 'b'}, {'c', 'd'}, {'d', 'e'}},
			kinds: []token.Kind{
				token.KindIdentifier, //reg_scalar
				token.KindLeftBrace,
				token.KindRightBrace,
			},
		},
		{
			keepFirstToken: true,
			name:           internal.CaseName("text_field", true, "message_list"),
			expectedObj: ast.TextField{ID: 12, Name: ast.Identifier{ID: 0}, Value: ast.TextMessageList{
				ID: 11, Values: []ast.TextMessage{
					{ID: 9},
					{ID: 10},
				},
			}},

			content: "messages [{}, {}]",
			indices: "a-------bcdefghijk",
			locs: [][2]rune{
				{'a', 'b'}, {'c', 'd'}, {'d', 'e'}, {'e', 'f'},
				{'f', 'g'}, {'h', 'i'}, {'i', 'j'}, {'j', 'k'},
			},
			kinds: []token.Kind{
				token.KindIdentifier, //messages
				token.KindLeftSquare,
				token.KindLeftBrace,
				token.KindRightBrace,
				token.KindComma,
				token.KindLeftBrace,
				token.KindRightBrace,
				token.KindRightSquare,
			},
		},
		{
			keepFirstToken: true,
			name:           internal.CaseName("text_field", false, "expected_colon"),
			expectedErrs: []error{
				gotUnexpected(&token.Token{ID: 1, Kind: token.KindInt}, token.KindColon),
			},

			content: "reg_scalar 10",
			indices: "a---------bc-d",
			locs:    [][2]rune{{'a', 'b'}, {'c', 'd'}},
			kinds: []token.Kind{
				token.KindIdentifier, //reg_scalar
				token.KindInt,
			},
		},
		{
			keepFirstToken: true,
			name:           internal.CaseName("text_field", false, "expected_colon_list"),
			expectedErrs: []error{
				gotUnexpected(&token.Token{ID: 1, Kind: token.KindLeftSquare}, token.KindColon),
			},

			content: "reg_scalar [10]",
			indices: "a---------bcd-ef",
			locs:    [][2]rune{{'a', 'b'}, {'c', 'd'}, {'d', 'e'}, {'e', 'f'}},
			kinds: []token.Kind{
				token.KindIdentifier, //reg_scalar
				token.KindLeftSquare,
				token.KindInt,
				token.KindRightSquare,
			},
		},
	}

	wrap := func(p *impl) (ast.TextField, error) { return p.parseTextField(1) }
	runTestCases(t, tests, checkTextField, wrap)
}

func TestParseTextMessage(t *testing.T) {
	tests := []TestCase[ast.TextMessage]{
		{
			keepFirstToken: true,
			name:           internal.CaseName("text_message", true, "brace"),
			expectedObj: ast.TextMessage{ID: 7, Fields: []ast.TextField{
				{ID: 6, Name: ast.Identifier{ID: 1}, Value: ast.Integer{ID: 3}},
			}},

			content: "{ reg_scalar: 10 }",
			indices: "abc---------def-ghi",
			locs:    [][2]rune{{'a', 'b'}, {'c', 'd'}, {'d', 'e'}, {'f', 'g'}, {'h', 'i'}},
			kinds: []token.Kind{
				token.KindLeftBrace,
				token.KindIdentifier, //reg_scalar
				token.KindColon,
				token.KindInt,
				token.KindRightBrace,
			},
		},
		{
			keepFirstToken: true,
			name:           internal.CaseName("text_message", true, "angle"),
			expectedObj: ast.TextMessage{ID: 7, Fields: []ast.TextField{
				{ID: 6, Name: ast.Identifier{ID: 1}, Value: ast.Integer{ID: 3}},
			}},

			content: "< reg_scalar: 10 >",
			indices: "abc---------def-ghi",
			locs:    [][2]rune{{'a', 'b'}, {'c', 'd'}, {'d', 'e'}, {'f', 'g'}, {'h', 'i'}},
			kinds: []token.Kind{
				token.KindLeftAngle,
				token.KindIdentifier, //reg_scalar
				token.KindColon,
				token.KindInt,
				token.KindRightAngle,
			},
		},
	}

	wrap := func(p *impl) (ast.TextMessage, error) { return p.parseTextMessage(1) }
	runTestCases(t, tests, checkTextMessage, wrap)
}
