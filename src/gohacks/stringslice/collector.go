package stringslice

// Collector accumulates string instances one at a time.
type Collector struct {
	data []string
}

// Add appends a string instance to this collector.
func (c *Collector) Add(text string) {
	c.data = append(c.data, text)
}

// Result provides all accumulated string instances.
func (c *Collector) Result() []string {
	if c.data == nil {
		return []string{}
	}
	return c.data
}
