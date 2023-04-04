package execute

import "fmt"

type Statistics interface {
	RegisterRun()
	PrintAnalysis()
}

func newStatistics(debug bool) Statistics {
	if debug {
		return &CommandsStatistics{}
	}
	return &NoStatistics{}
}

// CommandsStatistics is a Statistics implementation that counts how many commands were run.
type CommandsStatistics struct {
	CommandsCount int
}

func (s *CommandsStatistics) RegisterRun() {
	s.CommandsCount += 1
}

func (s *CommandsStatistics) PrintAnalysis() {
	fmt.Printf("Ran %d shell commands.", s.CommandsCount)
}

// NoStatistics is a statistics implementation that does nothing.
type NoStatistics struct{}

func (s *NoStatistics) RegisterRun() {}

func (s *NoStatistics) PrintAnalysis() {}
