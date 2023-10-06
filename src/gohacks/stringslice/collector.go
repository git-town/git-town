package stringslice

type Collector struct {
	data []string
}

func (c *Collector) Add(text string) {
	c.data = append(c.data, text)
}

func (c *Collector) Result() []string {
	if c.data == nil {
		return []string{}
	}
	return c.data
}
