package lexer

// TokenType is an alias type which tells of which kind the token is
type TokenType int

// These are all the token types
const (
	EOF          TokenType = iota - 1 // End Of File
	TokenIllegal                      // Illegal token
	TokenError                        // Error
	TokenSpace                        // Space (whitespace, '\n', '\r', '\t')
	TokenComment                      // Comment (single line or multiline)

	TokenIdentifier // Identifier
	TokenInt        // Integer
	TokenFloat      // Float
	TokenStr        // String ('...' or "...")

	TokenUnderscore  // _
	TokenEqual       // =
	TokenColon       // ,
	TokenSemicolon   // ;
	TokenDot         // .
	TokenLeftBrace   // {
	TokenRightBrace  // }
	TokenLeftSquare  // [
	TokenRightSquare // ]
	TokenLeftParen   // (
	TokenRightParen  // )
	TokenLeftAngle   // <
	TokenRightAngle  // >
)

var tokenTypeStr = [...]string{
	"EOF",
	"Illegal",
	"Error",
	"Space",
	"Comment",
	"Identifier",
	"Integer",
	"Float",
	"String",
	"_",
	"=",
	",",
	";",
	".",
	"{", "}",
	"[", "]",
	"(", ")",
	"<", ">",
}

func (t TokenType) String() string {
	return tokenTypeStr[t+1] // +1 because we start at iota - 1
}

// Token is a piece of the input
type Token struct {
	Type    TokenType
	Literal string
	Position
}
