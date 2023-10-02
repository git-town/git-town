package statistics

// None is a statistics implementation that does nothing.
type None struct{}

func (n *None) RegisterMessage(_ string) {}

func (n *None) RegisterRun() {}

func (n *None) PrintAnalysis() {}
