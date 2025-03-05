package parser

import (
	"fmt"
	"io"
	"iter"
	"math"
	"strings"

	"github.com/Clement-Jean/protein/lexer"
	"github.com/Clement-Jean/protein/source"
)

//go:generate stringer -type=NodeKind -linecomment
type NodeKind uint16

const (
	NodeKindUndefined    NodeKind = iota // Undefined
	NodeKindImportStmt                   // ImportStmt
	NodeKindPackageStmt                  // PackageStmt
	NodeKindMessageClose                 // MessageClose

	// type defs
	NodeKindMessageDecl      // MessageDecl
	NodeKindMessageOneOfDecl // OneOfDecl
	NodeKindEnumDecl         // EnumDecl
	NodeKindServiceDecl      // ServiceDecl

	// type refs
	NodeKindMessageFieldDecl // FieldDecl
	NodeKindEnumValueDecl    // ValueDecl
	NodeKindMapValue         // MapValue
	NodeKindRPCInputOutput   // RPCInputOutput
)

func (k NodeKind) IsTypeDef() bool {
	return k >= NodeKindMessageDecl && k < NodeKindMessageFieldDecl
}

func (k NodeKind) IsTypeRef() bool {
	return k >= NodeKindMessageFieldDecl
}

func (k NodeKind) NotType() bool {
	return k == NodeKindMessageOneOfDecl
}

type Node struct {
	TokIdx      uint32
	SubtreeSize uint32
	Kind        NodeKind
	HasError    bool
}

type ParseTree []Node

func (pt *ParseTree) roots() iter.Seq[int] {
	return func(yield func(int) bool) {
		i := len(*pt) - 1
		for i >= 0 {
			if !yield(i) {
				return
			}
			i -= int(math.Abs(float64((*pt)[i].SubtreeSize)))
		}
	}
}

func (pt *ParseTree) children(rootIdx int) iter.Seq[int] {
	return func(yield func(int) bool) {
		end := rootIdx - int((*pt)[rootIdx].SubtreeSize)
		i := rootIdx - 1
		for i > end {
			if !yield(i) {
				return
			}
			i -= int(math.Abs(float64((*pt)[i].SubtreeSize)))
		}
	}
}

func (pt *ParseTree) postorder() iter.Seq[int] {
	return func(yield func(int) bool) {
		for i := 0; i < len(*pt); i++ {
			if !yield(i) {
				return
			}
		}
	}
}

func (pt *ParseTree) printNode(out io.Writer, idx, depth int, toks *lexer.TokenizedBuffer, s *source.Buffer) bool {
	node := (*pt)[idx]
	indent := 2 * (depth + 1)

	fmt.Fprintf(out, "%s{", strings.Repeat(" ", indent))

	comment := ""

	if node.TokIdx > uint32(len(toks.TokenInfos)) {
		fmt.Fprint(out, "kind: <INSERT>")
	} else {
		kind := toks.TokenInfos[node.TokIdx].Kind

		if s != nil && (kind == lexer.TokenKindIdentifier || kind == lexer.TokenKindStr) {
			start := toks.TokenInfos[node.TokIdx].Offset
			end := toks.TokenInfos[node.TokIdx+1].Offset
			comment = fmt.Sprintf(" // %q", strings.TrimSpace(string(s.Range(start, end))))
		}

		fmt.Fprintf(out, "kind: %s", kind)
	}

	if node.HasError {
		fmt.Fprintf(out, ", hasError: %t", node.HasError)
	}

	if node.SubtreeSize > 1 {
		fmt.Fprintf(out, ", subtreeSize: %d", node.SubtreeSize)
	}

	fmt.Fprintf(out, "},%s\n", comment)
	return false
}

func (pt *ParseTree) Print(out io.Writer, s *source.Buffer, toks *lexer.TokenizedBuffer) {
	fmt.Fprintf(out, "parseTree = [\n")

	var stack []struct {
		tokIdx int
		depth  int
	}

	for node := range pt.roots() {
		stack = append(stack, struct {
			tokIdx int
			depth  int
		}{node, 0})
	}

	indents := make([]int, len(*pt))
	for len(stack) != 0 {
		top := stack[len(stack)-1]
		idx, depth := top.tokIdx, top.depth
		stack = stack[:len(stack)-1]

		for child := range pt.children(idx) {
			indents[child] = depth + 1
			stack = append(stack, struct {
				tokIdx int
				depth  int
			}{child, depth + 1})
		}
	}

	for node := range pt.postorder() {
		pt.printNode(out, node, indents[node], toks, s)
	}

	fmt.Fprintf(out, "]\n")
}
