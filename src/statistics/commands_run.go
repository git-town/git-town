package statistics

import "fmt"

// CommandsRun is a Statistics implementation that counts how many commands were run.
type CommandsRun struct {
	CommandsCount int
}

func (s *CommandsRun) RegisterRun() {
	s.CommandsCount++
}

func (s *CommandsRun) PrintAnalysis() {
	fmt.Printf("Ran %d shell commands.", s.CommandsCount)
}
