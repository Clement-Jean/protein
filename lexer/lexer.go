package lexer

// Lexer is protein's tokenizer
type Lexer interface {
	// NextToken returns the following token in the input source.
	NextToken() Token
}
