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

func (s *CommandsRun) RegisterRun() {
	s.CommandsCount++
}

func (s *CommandsRun) PrintAnalysis() {
	fmt.Printf("Ran %d shell commands.", s.CommandsCount)
}
