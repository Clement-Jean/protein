package lexer

//go:generate stringer -type=TokenKind -linecomment
type TokenKind uint8

const (
	TokenKindEOF     TokenKind = iota // EOF
	TokenKindBOF                      // BOF
	TokenKindError                    // Error
	TokenKindComment                  // Comment

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
	TokenKindStr   // String
)

// every token kind with the MSB set is also an identifier
const (
	TokenKindIdentifier   TokenKind = 128 + iota // Identifier
	TokenKindSyntax                              // syntax
	TokenKindEdition                             // edition
	TokenKindPackage                             // package
	TokenKindImport                              // import
	TokenKindPublic                              // public
	TokenKindWeak                                // weak
	TokenKindOption                              // option
	TokenKindReserved                            // reserved
	TokenKindMax                                 // max
	TokenKindEnum                                // enum
	TokenKindMessage                             // message
	TokenKindMap                                 // map
	TokenKindOneOf                               // oneof
	TokenKindExtensions                          // extensions
	TokenKindService                             // service
	TokenKindRpc                                 // rpc
	TokenKindReturns                             // returns
	TokenKindExtend                              // extend
	TokenKindTrue                                // true
	TokenKindFalse                               // false
	TokenKindTypeFloat                           // float
	TokenKindTypeDouble                          // double
	TokenKindTypeInt32                           // int32
	TokenKindTypeInt64                           // int64
	TokenKindTypeUint32                          // uint32
	TokenKindTypeUint64                          // uint64
	TokenKindTypeSint32                          // sint32
	TokenKindTypeSint64                          // sint64
	TokenKindTypeFixed32                         // fixed32
	TokenKindTypeFixed64                         // fixed64
	TokenKindTypeSfixed32                        // sfixed32
	TokenKindTypeSfixed64                        // sfixed64
	TokenKindTypeBool                            // bool
	TokenKindTypeString                          // string
	TokenKindTypeBytes                           // bytes
)

func (k TokenKind) IsIdentifier() bool {
	return k > 127
}
