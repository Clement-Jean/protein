package parser_test

import "testing"

func TestParseComment(t *testing.T) {
	tests := parseTestContent(t, "comment.txt")
	runParseTestCase(t, tests)
}
