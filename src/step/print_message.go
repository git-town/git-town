package step

import "fmt"

type PrintMessage struct {
	Message string
	Empty
}

func (step *PrintMessage) Run(_ RunArgs) error {
	fmt.Println(step.Message)
	return nil
}
