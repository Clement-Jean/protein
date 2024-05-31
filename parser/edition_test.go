package parser_test

import "testing"

func TestParseEdition(t *testing.T) {
	tests := parseTestContent(t, "edition.txt")
	runParseTestCase(t, tests)
}
