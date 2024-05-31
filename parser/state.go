package parser

import "fmt"

//go:generate stringer -type=state
type state uint8

const (
	stateTopLevel state = iota
	stateSyntaxAssign
	stateSyntaxFinish
	stateEditionAssign
	stateEditionFinish
	stateImportValue
	stateImportFinish
	statePackageFinish
	stateOptionAssign
	stateOptionFinish

	stateFullIdentifierRoot
	stateFullIdentifierRest
)

type stateStackEntry struct {
	st           state
	hasError     bool
	tokIdx       int
	subtreeStart int32
}

func (s stateStackEntry) String() string {
	return fmt.Sprintf(
		"{st: %s, hasError: %t, tokIdx: %d, subtreeStart: %d}",
		s.st, s.hasError, s.tokIdx, s.subtreeStart,
	)
}
