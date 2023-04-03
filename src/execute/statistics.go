package execute

import "fmt"

type Statistics struct {
	CommandsCount int
}

func (s *Statistics) RegisterRun() {
	if s != nil {
		s.CommandsCount += 1
	}
}

func (s *Statistics) PrintAnalysis() {
	fmt.Printf("Ran %d shell commands.", s.CommandsCount)
}

// NoStatistics is a mock statistics implementation for situations where no statistics are needed.
type NoStatistics struct{}

func (s *NoStatistics) RegisterRun() {}

func (s *NoStatistics) PrintAnalysis() {}
