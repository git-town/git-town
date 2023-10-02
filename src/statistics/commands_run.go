package statistics

import "fmt"

// CommandsRun is a Statistics implementation that counts how many commands were run.
type CommandsRun struct {
	CommandsCount int
	Messages      []string
}

func (s *CommandsRun) RegisterMessage(message string) {
	s.Messages = append(s.Messages, message)
}

func (cr *CommandsRun) RegisterRun() {
	cr.CommandsCount++
}

func (cr *CommandsRun) PrintAnalysis() {
	fmt.Printf("Ran %d shell commands.", cr.CommandsCount)
}
