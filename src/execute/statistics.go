package execute

import "fmt"

type Statistics interface {
	RegisterRun()
	PrintAnalysis()
}

// CommandsStatistics is a Statistics implementation that counts how many commands were run.
type CommandsStatistics struct {
	CommandsCount int
}

func (s *CommandsStatistics) RegisterRun() {
	s.CommandsCount++
}

func (s *CommandsStatistics) PrintAnalysis() {
	fmt.Printf("Ran %d shell commands.", s.CommandsCount)
}

// NoStatistics is a statistics implementation that does nothing.
type NoStatistics struct{}

func (s *NoStatistics) RegisterRun() {}

func (s *NoStatistics) PrintAnalysis() {}
