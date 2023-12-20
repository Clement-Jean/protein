package ast

import (
	"fmt"

	"github.com/Clement-Jean/protein/token"
)

type Rpc struct {
	Options        []Option
	Name           Identifier
	InputType      Identifier
	OutputType     Identifier
	ID             token.UniqueID
	IsServerStream bool
	IsClientStream bool
}

func (r Rpc) String() string {
	return fmt.Sprintf("{ type: Rpc, id: %d, name: %s, server_stream: %t, client_stream: %t }", r.ID, r.Name, r.IsServerStream, r.IsClientStream)
}

type Service struct {
	Options []Option
	Rpcs    []Rpc
	Name    Identifier
	ID      token.UniqueID
}

func (s Service) String() string {
	return fmt.Sprintf("{ type: Service, id: %d }", s.ID)
}
