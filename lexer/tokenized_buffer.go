package lexer

import "slices"

type LineIdx uint32

type TokenInfo struct {
	Offset uint32
	//  LineIdx LineIdx // LineInfo index inside TokenizedBuffer.LineInfos
	//  Column  uint32  // relative zero-based index from the beginning of a line
	Kind TokenKind
}

type LineInfo struct {
	Start uint32 // offset from the begining of the input text
	//  Len   uint32
}

type TokenizedBuffer struct {
	TokenInfos []TokenInfo
	LineInfos  []LineInfo
}

func (tb *TokenizedBuffer) FindLineIndex(offset uint32) LineIdx {
	idx, _ := slices.BinarySearchFunc(tb.LineInfos, offset, func(li LineInfo, offset uint32) int {
		if li.Start < offset {
			return -1
		} else if li.Start > offset {
			return 1
		}
		return 0
	})
	idx--
	return LineIdx(idx)
}

func (tb *TokenizedBuffer) GetIndentColumnNumber(idx LineIdx) uint32 {
	return tb.LineInfos[idx].Start + 1
}
