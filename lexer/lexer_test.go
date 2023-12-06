package lexer_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/Clement-Jean/protein/internal/span"
	"github.com/Clement-Jean/protein/lexer"
	"github.com/Clement-Jean/protein/token"
)

type TestCase struct {
	name  string
	input string
	kinds []token.Kind
	spans []span.Span
}

func runTestCases(t *testing.T, tests []TestCase) {
	t.Helper()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			l := lexer.New(strings.NewReader(test.input))
			kinds, spans := l.Tokenize()

			if !reflect.DeepEqual(test.kinds, kinds) {
				t.Fatalf("expected kinds: %v, got: %v", test.kinds, kinds)
			}

			if !reflect.DeepEqual(test.spans, spans) {
				t.Fatalf("expected spans: %v, got: %v", test.spans, spans)
			}
		})
	}
}

func symbolsTestCase() TestCase {
	const start = uint8(token.KindUnderscore)
	const end = uint8(token.KindSlash) + 1

	input := strings.Join(token.KindToStr[start:end], "")
	expectedTokens := make([]token.Kind, len(input)+1)
	expectedLocations := make([]span.Span, len(input)+1)

	for i := uint8(0); i < uint8(len(input)); i++ {
		expectedTokens[i] = token.Kind(start + i)
		expectedLocations[i] = span.Span{Start: uint64(i), End: uint64(i + 1)}
	}
	expectedTokens[len(input)] = token.KindEOF
	expectedLocations[len(input)] = span.Span{Start: uint64(len(input)), End: uint64(len(input))}

	return TestCase{"symbols", input, expectedTokens, expectedLocations}
}

func TestTokenize(t *testing.T) {
	tests := []TestCase{
		symbolsTestCase(),
		{
			name:  "illegal",
			input: "&",
			kinds: []token.Kind{token.KindIllegal, token.KindEOF},
			spans: []span.Span{{Start: 0, End: 1}, {Start: 1, End: 1}},
		},
		{
			name:  "spaces",
			input: "\t\n\v\f\n ",
			kinds: []token.Kind{token.KindSpace, token.KindEOF},
			spans: []span.Span{{Start: 0, End: 6}, {Start: 6, End: 6}},
		},
		{
			name:  "double_newline",
			input: "\n\n",
			kinds: []token.Kind{token.KindSpace, token.KindEOF},
			spans: []span.Span{{Start: 0, End: 2}, {Start: 2, End: 2}},
		},
		{
			name:  "line_comment_eof",
			input: "//this is a comment",
			kinds: []token.Kind{token.KindComment, token.KindEOF},
			spans: []span.Span{{Start: 0, End: 19}, {Start: 19, End: 19}},
		},
		{
			name:  "line_comment_newline",
			input: "//this is a comment\n",
			kinds: []token.Kind{token.KindComment, token.KindSpace, token.KindEOF},
			spans: []span.Span{{Start: 0, End: 19}, {Start: 19, End: 20}, {Start: 20, End: 20}},
		},
		{
			name:  "multiline_comment",
			input: "/*this is a comment*/",
			kinds: []token.Kind{token.KindComment, token.KindEOF},
			spans: []span.Span{{Start: 0, End: 21}, {Start: 21, End: 21}},
		},
		{
			name:  "multiline_comment_eof",
			input: "/*this is a comment",
			kinds: []token.Kind{token.KindErrorUnterminatedMultilineComment, token.KindEOF},
			spans: []span.Span{{Start: 0, End: 19}, {Start: 19, End: 19}},
		},
		{
			name:  "identifier",
			input: "hello_world2023 HelloWorld2023",
			kinds: []token.Kind{token.KindIdentifier, token.KindSpace, token.KindIdentifier, token.KindEOF},
			spans: []span.Span{{Start: 0, End: 15}, {Start: 15, End: 16}, {Start: 16, End: 30}, {Start: 30, End: 30}},
		},
		{
			name:  "string",
			input: "'test' \"test\"",
			kinds: []token.Kind{token.KindStr, token.KindSpace, token.KindStr, token.KindEOF},
			spans: []span.Span{{Start: 0, End: 6}, {Start: 6, End: 7}, {Start: 7, End: 13}, {Start: 13, End: 13}},
		},
		{
			name:  "escaped_string",
			input: "'this is a \\\"123string\\\"'",
			kinds: []token.Kind{token.KindStr, token.KindEOF},
			spans: []span.Span{{Start: 0, End: 25}, {Start: 25, End: 25}},
		},
		{
			name:  "unterminated_string_eof",
			input: "'test",
			kinds: []token.Kind{token.KindErrorUnterminatedQuotedString, token.KindEOF},
			spans: []span.Span{{Start: 0, End: 5}, {Start: 5, End: 5}},
		},
		{
			name:  "unterminated_string_newline",
			input: "'test\n'",
			kinds: []token.Kind{
				token.KindErrorUnterminatedQuotedString,
				token.KindSpace,
				token.KindErrorUnterminatedQuotedString,
				token.KindEOF,
			},
			spans: []span.Span{{Start: 0, End: 5}, {Start: 5, End: 6}, {Start: 6, End: 7}, {Start: 7, End: 7}},
		},
		{
			name:  "unterminated_string_mismatch",
			input: "\"test'",
			kinds: []token.Kind{token.KindErrorUnterminatedQuotedString, token.KindEOF},
			spans: []span.Span{{Start: 0, End: 6}, {Start: 6, End: 6}},
		},
		{
			name:  "decimal",
			input: "5 0 -5 +5",
			kinds: []token.Kind{
				token.KindInt,
				token.KindSpace,
				token.KindInt,
				token.KindSpace,
				token.KindInt,
				token.KindSpace,
				token.KindInt,
				token.KindEOF,
			},
			spans: []span.Span{
				{Start: 0, End: 1},
				{Start: 1, End: 2},
				{Start: 2, End: 3},
				{Start: 3, End: 4},
				{Start: 4, End: 6},
				{Start: 6, End: 7},
				{Start: 7, End: 9},
				{Start: 9, End: 9},
			},
		},
		{
			name:  "hexadecimal",
			input: "0xff 0XFF",
			kinds: []token.Kind{
				token.KindInt,
				token.KindSpace,
				token.KindInt,
				token.KindEOF,
			},
			spans: []span.Span{
				{Start: 0, End: 4},
				{Start: 4, End: 5},
				{Start: 5, End: 9},
				{Start: 9, End: 9},
			},
		},
		{
			name:  "octal",
			input: "056",
			kinds: []token.Kind{token.KindInt, token.KindEOF},
			spans: []span.Span{{Start: 0, End: 3}, {Start: 3, End: 3}},
		},
		{
			name:  "float",
			input: "-8.8 +0.8 -.8 +.8 .8 .8e8 .8e+8 .8e-8 8e8",
			kinds: []token.Kind{
				token.KindFloat,
				token.KindSpace,
				token.KindFloat,
				token.KindSpace,
				token.KindFloat,
				token.KindSpace,
				token.KindFloat,
				token.KindSpace,
				token.KindFloat,
				token.KindSpace,
				token.KindFloat,
				token.KindSpace,
				token.KindFloat,
				token.KindSpace,
				token.KindFloat,
				token.KindSpace,
				token.KindFloat,
				token.KindEOF,
			},
			spans: []span.Span{
				{Start: 0, End: 4},
				{Start: 4, End: 5},
				{Start: 5, End: 9},
				{Start: 9, End: 10},
				{Start: 10, End: 13},
				{Start: 13, End: 14},
				{Start: 14, End: 17},
				{Start: 17, End: 18},
				{Start: 18, End: 20},
				{Start: 20, End: 21},
				{Start: 21, End: 25},
				{Start: 25, End: 26},
				{Start: 26, End: 31},
				{Start: 31, End: 32},
				{Start: 32, End: 37},
				{Start: 37, End: 38},
				{Start: 38, End: 41},
				{Start: 41, End: 41},
			},
		},
	}

	runTestCases(t, tests)
}
