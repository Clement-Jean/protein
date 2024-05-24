package lexer

//go:generate stringer -type=TokenKind
type TokenKind uint8

const (
	TokenKindEOF     TokenKind = iota // End Of File
	TokenKindBOF                      // Begining of file
	TokenKindError                    // Error
	TokenKindComment                  // Comment (single line or multiline)

	TokenKindIdentifier // Identifier
	TokenKindInt        // Integer
	TokenKindFloat      // Float
	TokenKindStr        // String ('...' or "...")

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
)
