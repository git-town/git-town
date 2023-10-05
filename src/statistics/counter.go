package statistics

// Counter is a Statistics implementation that counts how many commands were run.
type Counter struct {
	count int
}

func (cr *Counter) Count() int {
	return cr.count
}

func (cr *Counter) RegisterRun() {
	cr.count++
}
