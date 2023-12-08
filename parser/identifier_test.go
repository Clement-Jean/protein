package parser

import (
	"github.com/Clement-Jean/protein/ast"
	"log"
	"slices"
	"testing"

	"github.com/Clement-Jean/protein/token"
)

func checkIdentifier(t *testing.T, a ast.Identifier, b ast.Identifier) {
	t.Helper()
	checkIDs(t, a.ID, b.ID)
	checkIdentifierParts(t, a.Parts, b.Parts)
}

func checkIdentifierParts(t *testing.T, parts []token.UniqueID, expected []token.UniqueID) {
	t.Helper()
	if len(parts) != len(expected) {
		for _, part := range parts {
			log.Println(part)
		}
		t.Fatalf("expected %d parts, got '%d'", len(expected), len(parts))
	}

	slices.Sort(parts)
	slices.Sort(expected)
	for i, part := range parts {
		checkIDs(t, part, expected[i])
	}
}
