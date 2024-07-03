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
	TokenKindLeftSquare  // [
	TokenKindLeftParen   // (
	TokenKindLeftAngle   // <
	TokenKindRightBrace  // }
	TokenKindRightSquare // ]
	TokenKindRightParen  // )
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
	TokenKindTo                                  // to
	TokenKindMax                                 // max
	TokenKindEnum                                // enum
	TokenKindMessage                             // message
	TokenKindMap                                 // map
	TokenKindRepeated                            // repeated
	TokenKindOptional                            // optional
	TokenKindOneOf                               // oneof
	TokenKindExtensions                          // extensions
	TokenKindService                             // service
	TokenKindRPC                                 // rpc
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

var literals = []string{
	"syntax",
	"edition",
	"package",
	"import",
	"public",
	"weak",
	"option",
	"reserved",
	"to",
	"max",
	"enum",
	"message",
	"map",
	"repeated",
	"optional",
	"oneof",
	"extensions",
	"service",
	"rpc",
	"returns",
	"extend",
	"true",
	"false",
	"float",
	"double",
	"int32",
	"int64",
	"uint32",
	"uint64",
	"sint32",
	"sint64",
	"fixed32",
	"fixed64",
	"sfixed32",
	"sfixed64",
	"bool",
	"string",
	"bytes",
}
var kinds = [...]TokenKind{
	TokenKindSyntax,
	TokenKindEdition,
	TokenKindPackage,
	TokenKindImport,
	TokenKindPublic,
	TokenKindWeak,
	TokenKindOption,
	TokenKindReserved,
	TokenKindTo,
	TokenKindMax,
	TokenKindEnum,
	TokenKindMessage,
	TokenKindMap,
	TokenKindRepeated,
	TokenKindOptional,
	TokenKindOneOf,
	TokenKindExtensions,
	TokenKindService,
	TokenKindRPC,
	TokenKindReturns,
	TokenKindExtend,
	TokenKindTrue,
	TokenKindFalse,
	TokenKindTypeFloat,
	TokenKindTypeDouble,
	TokenKindTypeInt32,
	TokenKindTypeInt64,
	TokenKindTypeUint32,
	TokenKindTypeUint64,
	TokenKindTypeSint32,
	TokenKindTypeSint64,
	TokenKindTypeFixed32,
	TokenKindTypeFixed64,
	TokenKindTypeSfixed32,
	TokenKindTypeSfixed64,
	TokenKindTypeBool,
	TokenKindTypeString,
	TokenKindTypeBytes,
}

func (k TokenKind) IsIdentifier() bool {
	return k >= TokenKindIdentifier
}

func (k TokenKind) IsOpeningSymbol() bool {
	return k >= TokenKindLeftBrace && k <= TokenKindLeftAngle
}

func (k TokenKind) IsClosingSymbol() bool {
	return k >= TokenKindRightBrace && k <= TokenKindRightAngle
}

func (k TokenKind) MatchingClosingSymbol() TokenKind {
	mid := (TokenKindRightAngle-TokenKindLeftBrace)/2 + 1
	pos := k - TokenKindLeftBrace
	return TokenKindLeftBrace + mid + pos
}
