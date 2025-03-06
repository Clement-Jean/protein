package lexer

import "fmt"

type InvalidChar struct {
	Character byte
}

func (ic *InvalidChar) Error() string {
	return fmt.Sprintf("invalid character %c", ic.Character)
}
