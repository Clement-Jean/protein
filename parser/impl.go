package parser

import (
	"fmt"
	"os"
	"strings"

	"github.com/Clement-Jean/protein/lexer"
	"golang.org/x/exp/slices"
	pb "google.golang.org/protobuf/types/descriptorpb"
)

// Impl is the implementation for the Parser interface.
type Impl struct {
	l         lexer.Lexer
	curToken  lexer.Token
	peekToken lexer.Token
}

// New creates a new instance of the Parser
func New(l lexer.Lexer) Parser {
	p := &Impl{l: l}
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Impl) nextToken() {
	p.curToken = p.peekToken

	for p.curToken.Type == lexer.TokenSpace {
		p.curToken = p.l.NextToken()
	}

	p.peekToken = p.l.NextToken()

	for p.peekToken.Type == lexer.TokenSpace {
		p.peekToken = p.l.NextToken()
	}
}

func (p *Impl) accept(original lexer.TokenType, expected ...lexer.TokenType) bool {
	if !slices.Contains(expected, original) {
		p.error(fmt.Sprintf(
			errorUnexpectedPeek,
			strings.Trim(fmt.Sprint(expected), "[]"),
			original.String(),
		))
		return false
	}

	p.nextToken()
	return true
}

// acceptPeek returns true and advance token
// if the type t is equal to the peekToken.Type
// else it returns false
func (p *Impl) acceptPeek(tt ...lexer.TokenType) bool {
	return p.accept(p.peekToken.Type, tt...)
}

func (p *Impl) error(msg string) {
	fmt.Fprint(os.Stderr, msg)
}

var parseFuncs = map[string]func(p *Impl, d *pb.FileDescriptorProto){
	"syntax":  func(p *Impl, d *pb.FileDescriptorProto) { d.Syntax = p.parseSyntax() },
	"package": func(p *Impl, d *pb.FileDescriptorProto) { d.Package = p.parsePackage() },
	"import": func(p *Impl, d *pb.FileDescriptorProto) {
		dep := p.parseImport()
		if len(dep) != 0 {
			d.Dependency = append(d.Dependency, dep)
		}
	},
}

// Parse populates a FileDescriptorProto
func (p *Impl) Parse() pb.FileDescriptorProto {
	d := pb.FileDescriptorProto{}

	for p.curToken.Type != lexer.EOF {
		if p.curToken.Type == lexer.TokenIdentifier {
			fn, ok := parseFuncs[p.curToken.Literal]
			if !ok {
				//TODO: recovering
				//p.error(fmt.Sprintf(errorUnknownKeyword, p.peekToken.Literal))
				break
			}
			fn(p, &d)
		}
		p.nextToken()
	}

	return d
}
