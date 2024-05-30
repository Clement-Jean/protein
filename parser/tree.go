package parser

type Node struct {
	TokIdx      int
	SubtreeSize int32
	HasError    bool
}

type ParseTree []Node
