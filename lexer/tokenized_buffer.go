package lexer

type LineIdx uint32

type TokenInfo struct {
	LineIdx LineIdx // LineInfo index inside TokenizedBuffer.LineInfos
	Column  uint32  // relative zero-based index from the beginning of a line
	Kind    TokenKind
}

type LineInfo struct {
	Start int // offset from the begining of the input text
	Len   uint32
}

type TokenizedBuffer struct {
	TokenInfos []TokenInfo
	LineInfos  []LineInfo
}
