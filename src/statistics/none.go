package statistics

// None is a statistics implementation that does nothing.
type None struct{}

func (s *None) RegisterRun() {}

func (s *None) PrintAnalysis() {}
