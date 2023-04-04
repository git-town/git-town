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

type CommandsStatistics struct {
	CommandsCount int
}

func (s *CommandsStatistics) RegisterRun() {
	if s != nil {
		s.CommandsCount += 1
	}
}

func (s *CommandsStatistics) PrintAnalysis() {
	fmt.Printf("Ran %d shell commands.", s.CommandsCount)
}

// NoStatistics is a mock statistics implementation for situations where no statistics are needed.
type NoStatistics struct{}

func (s *NoStatistics) RegisterRun() {}

func (s *NoStatistics) PrintAnalysis() {}
