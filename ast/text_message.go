package ast

import (
	"fmt"

	"github.com/Clement-Jean/protein/token"
)

type TextField struct {
	Value Expression
	Name  Identifier
	ID    token.UniqueID
}

func (f TextField) declarationNode()      {}
func (f TextField) GetID() token.UniqueID { return f.ID }
func (f TextField) String() string {
	return fmt.Sprintf("{ type: TextField, id: %d }", f.ID)
}

type TextScalarList struct {
	Values []Expression
	ID     token.UniqueID
}

func (l TextScalarList) expressionNode()       {}
func (l TextScalarList) GetID() token.UniqueID { return l.ID }
func (l TextScalarList) String() string {
	return fmt.Sprintf("{ type: ScalarList, id: %d }", l.ID)
}

type TextMessageList struct {
	Values []TextMessage
	ID     token.UniqueID
}

func (l TextMessageList) expressionNode()       {}
func (l TextMessageList) GetID() token.UniqueID { return l.ID }
func (l TextMessageList) String() string {
	return fmt.Sprintf("{ type: MessageList, id: %d }", l.ID)
}

type TextMessage struct {
	Fields []TextField
	ID     token.UniqueID
}

func (m TextMessage) expressionNode()       {}
func (m TextMessage) GetID() token.UniqueID { return m.ID }
func (m TextMessage) String() string {
	return fmt.Sprintf("{ type: TextMessage, id: %d }", m.ID)
}
