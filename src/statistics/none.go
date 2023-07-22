package statistics

// NoStatistics is a statistics implementation that does nothing.
type NoStatistics struct{}

func (s *NoStatistics) RegisterRun() {}

func (s *NoStatistics) PrintAnalysis() {}
