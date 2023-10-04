package steps

import "fmt"

type PrintMessageStep struct {
	Message string
	EmptyStep
}

func (step *PrintMessageStep) Run(_ RunArgs) error {
	fmt.Println(step.Message)
	return nil
}
