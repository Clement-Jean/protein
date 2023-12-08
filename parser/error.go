package parser

import (
	"strings"

	"github.com/Clement-Jean/protein/token"
)

type Error struct {
	Msg string
	ID  token.UniqueID
}

func (e *Error) Error() string {
	return e.Msg
}

func gotUnexpected(got *token.Token, expected ...token.Kind) error {
	if len(expected) == 0 {
		return nil
	}

	var sb strings.Builder

	sb.WriteString("expected ")

	if len(expected) > 1 {
		sb.WriteString("[ ")

		for _, e := range expected {
			sb.WriteRune('\'')
			sb.WriteString(e.String())
			sb.WriteString("' ")
		}

		sb.WriteString("], ")
	} else {
		sb.WriteRune('\'')
		sb.WriteString(expected[0].String())
		sb.WriteString("', ")
	}

	sb.WriteString("got '")
	sb.WriteString(got.Kind.String())
	sb.WriteRune('\'')
	return &Error{ID: got.ID, Msg: sb.String()}
}
