package parser

import "fmt"

//go:generate stringer -type=state
type state uint8

const (
	stateTopLevel state = iota

	// SYNTAXES
	stateSyntaxAssign
	stateSyntaxFinish

	// EDITIONS
	stateEditionAssign
	stateEditionFinish

	// IMPORTS
	stateImportValue
	stateImportFinish

	// PACKAGES
	statePackageFinish

	// OPTIONS
	stateOptionName
	stateOptionNameRest
	stateOptionNameParenFinish
	stateOptionAssign
	stateOptionEqual
	stateOptionFinish
	stateTextFieldValue
	stateTextFieldAssign
	stateTextFieldName
	stateTextFieldExtensionName
	stateTextFieldExtensionNameFinish
	stateTextMessageValue
	stateTextMessageInsertSemicolon
	stateTextMessageFinishRightBrace
	stateTextMessageFinishRightAngle
	stateTextListValue
	stateTextListFinish

	// MESSAGES
	stateMessageBlock
	stateMessageFieldAssign
	stateMessageFieldOption
	stateMessageFieldOptionAssign
	stateMessageFieldOptionFinish
	stateMessageFieldFinish
	stateMessageMapStart
	stateMessageMapKeyValue
	stateMessageMapComma
	stateMessageMapFinish
	stateMessageValue
	stateMessageFinish

	// RESERVEDS
	stateReservedRange
	stateReservedName
	stateReservedFinish

	// ONEOFS
	stateOneofBlock
	stateOneofValue
	stateOneofFinish

	// ENUMS
	stateEnumBlock
	stateEnumValue
	stateEnumFinish

	// SERVICES
	stateServiceBlock
	stateServiceValue
	stateServiceFinish

	// RPCS
	stateRPCDefinition
	stateRPCReqRes
	stateRPCReqResFinish
	stateRPCValue
	stateRPCFinish

	// IDENTIFIERS
	stateIdentifier
	stateFullIdentifierRoot
	stateFullIdentifierRest

	// MISC
	stateEnder
)

type stateStackEntry struct {
	st           state
	hasError     bool
	kind         NodeKind
	tokIdx       uint32
	subtreeStart uint32
}

func (s stateStackEntry) String() string {
	return fmt.Sprintf(
		"{st: %s, hasError: %t, tokIdx: %d, subtreeStart: %d}",
		s.st, s.hasError, s.tokIdx, s.subtreeStart,
	)
}
