package math

import (
	"fmt"
	"turboengine/common/protocol"
)

type Math struct {
	Ver   string `version:"1.0"`
	Do    func(src protocol.Mailbox, dest protocol.Mailbox, x int, y int) (int, error)
	Print func(src protocol.Mailbox, dest protocol.Mailbox, str string) error
	XXX   interface{}
}

type MathService struct {
}

func (m *MathService) Do(src protocol.Mailbox, dest protocol.Mailbox, x int, y int) (int, error) {
	return x + y, nil
}

func (m *MathService) Print(src protocol.Mailbox, dest protocol.Mailbox, str string) error {
	fmt.Println(str)
	return nil
}
