package gohacks

// Counter is a Statistics implementation that counts how many commands were run.
type Counter struct {
	count int
}

func (c *Counter) Count() int {
	return c.count
}

func (c *Counter) Register() {
	c.count++
}
