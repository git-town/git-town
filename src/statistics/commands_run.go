package statistics

import "fmt"

// CommandsRun is a Statistics implementation that counts how many commands were run.
type CommandsRun struct {
	CommandsCount int
	messages      []string
}

func (cr *CommandsRun) RegisterMessage(message string) {
	cr.messages = append(cr.messages, message)
}

func (cr *CommandsRun) RegisterRun() {
	cr.CommandsCount++
}

func (cr *CommandsRun) PrintAnalysis() {
	fmt.Printf("Ran %d shell commands.", cr.CommandsCount)
}

func (cr *CommandsRun) PrintMessages() {
	for _, message := range cr.messages {
		fmt.Println(message)
	}
}
