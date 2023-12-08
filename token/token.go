package token

// Kind is an alias type which tells of which kind the token is.
type Kind uint8

const (
	KindEOF                           Kind = iota // End Of File
	KindIllegal                                   // Illegal token
	KindErrorUnterminatedQuotedString             // Error
	KindErrorUnterminatedMultilineComment
	KindSpace   // Space (whitespace, '\n', '\r', '\t')
	KindComment // Comment (single line or multiline)

	KindIdentifier // Identifier
	KindInt        // Integer
	KindFloat      // Float
	KindStr        // String ('...' or "...")

	KindUnderscore  // _
	KindEqual       // =
	KindComma       // ,
	KindColon       // :
	KindSemicolon   // ;
	KindDot         // .
	KindLeftBrace   // {
	KindRightBrace  // }
	KindLeftSquare  // [
	KindRightSquare // ]
	KindLeftParen   // (
	KindRightParen  // )
	KindLeftAngle   // <
	KindRightAngle  // >
	KindSlash

	KindSyntax
)

var KindToStr = [...]string{
	"EOF",
	"Illegal",
	"Unterminated quoted string Error",
	"Unterminated multiline comment Error",
	"Space",
	"Comment",
	"Identifier",
	"Int",
	"Float",
	"String",
	"_",
	"=",
	",",
	":",
	";",
	".",
	"{",
	"}",
	"[",
	"]",
	"(",
	")",
	"<",
	">",
	"/",
	"syntax",
}

func (k Kind) IsSymbol() bool { return KindUnderscore < k && k < KindSlash }
func (k Kind) String() string { return KindToStr[k] }

type UniqueID = int

type Token struct {
	ID   UniqueID
	Kind Kind
}
