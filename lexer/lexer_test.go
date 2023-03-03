package lexer

import (
	"strings"
	"testing"
)

type Check struct {
	expectedType     TokenType
	expectedLiteral  string
	expectedPosition Position
}

func runChecks(t *testing.T, l Lexer, tests []Check) {
	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected='%s', got='%s'", i, tt.expectedLiteral, tok.Literal)
		}

		if tok.Position.Offset != tt.expectedPosition.Offset {
			t.Fatalf("tests[%d] - offset wrong. expected=%d, got=%d", i, tt.expectedPosition.Offset, tok.Position.Offset)
		}

		if tok.Position.Line != tt.expectedPosition.Line {
			t.Fatalf("tests[%d] - line wrong. expected=%d, got=%d", i, tt.expectedPosition.Line, tok.Position.Line)
		}

		if tok.Position.Column != tt.expectedPosition.Column {
			t.Fatalf("tests[%d] - column wrong. expected=%d, got=%d", i, tt.expectedPosition.Column, tok.Position.Column)
		}
	}
}

func TestNextTokenOnSymbols(t *testing.T) {
	input := strings.Join(tokenTypeStr[TokenUnderscore+1:TokenRightAngle+2], "")
	runChecks(t, New(input), []Check{
		{TokenUnderscore, "_", Position{0, 1, 0}},
		{TokenEqual, "=", Position{1, 1, 1}},
		{TokenColon, ",", Position{2, 1, 2}},
		{TokenSemicolon, ";", Position{3, 1, 3}},
		{TokenDot, ".", Position{4, 1, 4}},
		{TokenLeftBrace, "{", Position{5, 1, 5}},
		{TokenRightBrace, "}", Position{6, 1, 6}},
		{TokenLeftSquare, "[", Position{7, 1, 7}},
		{TokenRightSquare, "]", Position{8, 1, 8}},
		{TokenLeftParen, "(", Position{9, 1, 9}},
		{TokenRightParen, ")", Position{10, 1, 10}},
		{TokenLeftAngle, "<", Position{11, 1, 11}},
		{TokenRightAngle, ">", Position{12, 1, 12}},
		{EOF, "", Position{13, 1, 13}},
	})
}

func TestNextTokenOnSpace(t *testing.T) {
	runChecks(t, New("\t\n\v\f\r "), []Check{
		{TokenSpace, "\t\n\v\f\r ", Position{0, 1, 0}},
		{EOF, "", Position{6, 2, 4}},
	})
}

func TestNextTokenOnMultipleNewline(t *testing.T) {
	runChecks(t, New("\n\n"), []Check{
		{TokenSpace, "\n\n", Position{0, 1, 0}},
		{EOF, "", Position{2, 3, 0}},
	})
}

func TestNextTokenOnLineCommentWithEOF(t *testing.T) {
	runChecks(t, New("//this is a comment"), []Check{
		{TokenComment, "//this is a comment", Position{0, 1, 0}},
		{EOF, "", Position{19, 1, 19}},
	})
}

func TestNextTokenOnLineCommentWithNewLine(t *testing.T) {
	runChecks(t, New("//this is a comment\n_"), []Check{
		{TokenComment, "//this is a comment", Position{0, 1, 0}},
		{TokenSpace, "\n", Position{19, 1, 19}},
		{TokenUnderscore, "_", Position{20, 2, 0}},
		{EOF, "", Position{21, 2, 1}},
	})
}

func TestNextTokenOnMultilineComment(t *testing.T) {
	runChecks(t, New("/*this is a comment*/_"), []Check{
		{TokenComment, "/*this is a comment*/", Position{0, 1, 0}},
		{TokenUnderscore, "_", Position{21, 1, 21}},
		{EOF, "", Position{22, 1, 22}},
	})
}

func TestNextTokenOnUnterminatedMultilineComment(t *testing.T) {
	runChecks(t, New("/*this is a comment"), []Check{
		{TokenError, errorUnterminatedMultilineComment, Position{0, 1, 0}},
	})
}

func TestNextTokenOnIdentifier(t *testing.T) {
	runChecks(t, New("hello_world2023 HelloWorld2023"), []Check{
		{TokenIdentifier, "hello_world2023", Position{0, 1, 0}},
		{TokenSpace, " ", Position{15, 1, 15}},
		{TokenIdentifier, "HelloWorld2023", Position{16, 1, 16}},
		{EOF, "", Position{30, 1, 30}},
	})
}

func TestNextTokenOnString(t *testing.T) {
	runChecks(t, New("'test' \"test\""), []Check{
		{TokenStr, "'test'", Position{0, 1, 0}},
		{TokenSpace, " ", Position{6, 1, 6}},
		{TokenStr, "\"test\"", Position{7, 1, 7}},
		{EOF, "", Position{13, 1, 13}},
	})
}

func TestNextTokenOnUnterminatedString(t *testing.T) {
	runChecks(t, New("'test"), []Check{
		{TokenError, errorUnterminatedQuotedString, Position{0, 1, 0}},
	})
}

func TestNextTokenOnMismatchedQuotesString(t *testing.T) {
	runChecks(t, New("\"test'"), []Check{
		{TokenError, errorUnterminatedQuotedString, Position{0, 1, 0}},
	})
}

func TestNextTokenOnIntDecimal(t *testing.T) {
	runChecks(t, New("5 0 -5 +5"), []Check{
		{TokenInt, "5", Position{0, 1, 0}},
		{TokenSpace, " ", Position{1, 1, 1}},
		{TokenInt, "0", Position{2, 1, 2}},
		{TokenSpace, " ", Position{3, 1, 3}},
		{TokenInt, "-5", Position{4, 1, 4}},
		{TokenSpace, " ", Position{6, 1, 6}},
		{TokenInt, "+5", Position{7, 1, 7}},
		{EOF, "", Position{9, 1, 9}},
	})
}

func TestNextTokenOnIntHex(t *testing.T) {
	runChecks(t, New("0xff 0XFF"), []Check{
		{TokenInt, "0xff", Position{0, 1, 0}},
		{TokenSpace, " ", Position{4, 1, 4}},
		{TokenInt, "0XFF", Position{5, 1, 5}},
		{EOF, "", Position{9, 1, 9}},
	})
}

func TestNextTokenOnIntOctal(t *testing.T) {
	runChecks(t, New("056"), []Check{
		{TokenInt, "056", Position{0, 1, 0}},
		{EOF, "", Position{3, 1, 3}},
	})
}

func TestNextTokenOnFloat(t *testing.T) {
	runChecks(t, New("-0.8 +0.8 -.8 +.8 .8 .8e8 .8e+8 .8e-8 8e8"), []Check{
		{TokenFloat, "-0.8", Position{0, 1, 0}},
		{TokenSpace, " ", Position{4, 1, 4}},
		{TokenFloat, "+0.8", Position{5, 1, 5}},
		{TokenSpace, " ", Position{9, 1, 9}},
		{TokenFloat, "-.8", Position{10, 1, 10}},
		{TokenSpace, " ", Position{13, 1, 13}},
		{TokenFloat, "+.8", Position{14, 1, 14}},
		{TokenSpace, " ", Position{17, 1, 17}},
		{TokenFloat, ".8", Position{18, 1, 18}},
		{TokenSpace, " ", Position{20, 1, 20}},
		{TokenFloat, ".8e8", Position{21, 1, 21}},
		{TokenSpace, " ", Position{25, 1, 25}},
		{TokenFloat, ".8e+8", Position{26, 1, 26}},
		{TokenSpace, " ", Position{31, 1, 31}},
		{TokenFloat, ".8e-8", Position{32, 1, 32}},
		{TokenSpace, " ", Position{37, 1, 37}},
		{TokenFloat, "8e8", Position{38, 1, 38}},
		{EOF, "", Position{41, 1, 41}},
	})
}
