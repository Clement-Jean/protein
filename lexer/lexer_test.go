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
	}

	runTestCases(t, tests)
}
