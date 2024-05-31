package parser_test

import "testing"

func TestParsePackage(t *testing.T) {
	tests := parseTestContent(t, "package.txt")
	runParseTestCase(t, tests)
}
