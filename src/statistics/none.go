package statistics

// None is a statistics implementation that does nothing.
type None struct{}

func (n *None) RegisterMessage(message string) {}

func (n *None) RegisterRun() {}

func (n *None) PrintAnalysis() {}
