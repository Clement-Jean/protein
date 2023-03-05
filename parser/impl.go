package parser

import (
	"github.com/Clement-Jean/protein/lexer"
	pb "google.golang.org/protobuf/types/descriptorpb"
)

// Impl is the implementation for the Parser interface.
type Impl struct {
	l         lexer.Lexer
	curToken  lexer.Token
	peekToken lexer.Token
}

func New(l lexer.Lexer) Parser {
	p := &Impl{l: l}
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Impl) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Impl) expectPeek(t lexer.TokenType) bool {
	if p.peekToken.Type == t {
		p.nextToken()
		return true
	} else {
		return false
	}
}

func (p *Impl) Parse() pb.FileDescriptorProto {
	d := pb.FileDescriptorProto{}

	for p.curToken.Type != lexer.EOF {
		if p.curToken.Type == lexer.TokenIdentifier {
			switch p.curToken.Literal {
			case "syntax":
				d.Syntax = p.parseSyntax()
			default:
				break
			}
		}
		p.nextToken()
	}

	return d
}
