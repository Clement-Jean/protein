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

		// asserts on Position for later
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
