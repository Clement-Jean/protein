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

// /!\ BEWARE: All the tokens and literals after this line
// are sorted in alphabetical order. This is important for
// the lexing algorithm (binary search)!

// every token kind with the MSB set is also an identifier
const (
	TokenKindIdentifier   TokenKind = 128 + iota // Identifier
	TokenKindTypeBool                            // bool
	TokenKindTypeBytes                           // bytes
	TokenKindTypeDouble                          // double
	TokenKindEdition                             // edition
	TokenKindEnum                                // enum
	TokenKindExtend                              // extend
	TokenKindExtensions                          // extensions
	TokenKindFalse                               // false
	TokenKindTypeFixed32                         // fixed32
	TokenKindTypeFixed64                         // fixed64
	TokenKindTypeFloat                           // float
	TokenKindImport                              // import
	TokenKindTypeInt32                           // int32
	TokenKindTypeInt64                           // int64
	TokenKindMap                                 // map
	TokenKindMax                                 // max
	TokenKindMessage                             // message
	TokenKindOneOf                               // oneof
	TokenKindOption                              // option
	TokenKindOptional                            // optional
	TokenKindPackage                             // package
	TokenKindPublic                              // public
	TokenKindRepeated                            // repeated
	TokenKindReserved                            // reserved
	TokenKindReturns                             // returns
	TokenKindRPC                                 // rpc
	TokenKindService                             // service
	TokenKindTypeSfixed32                        // sfixed32
	TokenKindTypeSfixed64                        // sfixed64
	TokenKindTypeSint32                          // sint32
	TokenKindTypeSint64                          // sint64
	TokenKindStream                              // stream
	TokenKindTypeString                          // string
	TokenKindSyntax                              // syntax
	TokenKindTo                                  // to
	TokenKindTrue                                // true
	TokenKindTypeUint32                          // uint32
	TokenKindTypeUint64                          // uint64
	TokenKindWeak                                // weak
)

var literals = []string{
	"bool",
	"bytes",
	"double",
	"edition",
	"enum",
	"extend",
	"extensions",
	"false",
	"fixed32",
	"fixed64",
	"float",
	"import",
	"int32",
	"int64",
	"map",
	"max",
	"message",
	"oneof",
	"option",
	"optional",
	"package",
	"public",
	"repeated",
	"reserved",
	"returns",
	"rpc",
	"service",
	"sfixed32",
	"sfixed64",
	"sint32",
	"sint64",
	"stream",
	"string",
	"syntax",
	"to",
	"true",
	"uint32",
	"uint64",
	"weak",
}
var kinds = [...]TokenKind{
	TokenKindTypeBool,
	TokenKindTypeBytes,
	TokenKindTypeDouble,
	TokenKindEdition,
	TokenKindEnum,
	TokenKindExtend,
	TokenKindExtensions,
	TokenKindFalse,
	TokenKindTypeFixed32,
	TokenKindTypeFixed64,
	TokenKindTypeFloat,
	TokenKindImport,
	TokenKindTypeInt32,
	TokenKindTypeInt64,
	TokenKindMap,
	TokenKindMax,
	TokenKindMessage,
	TokenKindOneOf,
	TokenKindOption,
	TokenKindOptional,
	TokenKindPackage,
	TokenKindPublic,
	TokenKindRepeated,
	TokenKindReserved,
	TokenKindReturns,
	TokenKindRPC,
	TokenKindService,
	TokenKindTypeSfixed32,
	TokenKindTypeSfixed64,
	TokenKindTypeSint32,
	TokenKindTypeSint64,
	TokenKindStream,
	TokenKindTypeString,
	TokenKindSyntax,
	TokenKindTo,
	TokenKindTrue,
	TokenKindTypeUint32,
	TokenKindTypeUint64,
	TokenKindWeak,
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
