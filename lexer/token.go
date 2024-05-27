package lexer

//go:generate stringer -type=TokenKind
type TokenKind uint8

const (
	TokenKindEOF     TokenKind = iota // End Of File
	TokenKindBOF                      // Begining of file
	TokenKindError                    // Error
	TokenKindComment                  // Comment (single line or multiline)

	TokenKindUnderscore  // _
	TokenKindEqual       // =
	TokenKindComma       // ,
	TokenKindColon       // :
	TokenKindSemicolon   // ;
	TokenKindDot         // .
	TokenKindLeftBrace   // {
	TokenKindRightBrace  // }
	TokenKindLeftSquare  // [
	TokenKindRightSquare // ]
	TokenKindLeftParen   // (
	TokenKindRightParen  // )
	TokenKindLeftAngle   // <
	TokenKindRightAngle  // >
	TokenKindSlash       // /

	TokenKindInt   // Integer
	TokenKindFloat // Float
	TokenKindStr   // String ('...' or "...")
)

// every token kind with the MSB set is also an identifier
const (
	TokenKindIdentifier TokenKind = 128 + iota // Identifier
	TokenKindSyntax
	TokenKindEdition
	TokenKindPackage
	TokenKindImport
	TokenKindPublic
	TokenKindWeak
	TokenKindOption
	TokenKindReserved
	TokenKindMax
	TokenKindEnum
	TokenKindMessage
	TokenKindMap
	TokenKindOneOf
	TokenKindExtensions
	TokenKindService
	TokenKindRpc
	TokenKindReturns
	TokenKindExtend
)

func (k TokenKind) isIdentifier() bool {
	return k > 127
}
