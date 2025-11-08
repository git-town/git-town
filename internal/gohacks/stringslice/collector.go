package stringslice

// Collector accumulates strings.
type Collector struct {
	data *[]string
}

func NewCollector() Collector {
	var data []string
	return Collector{
		data: &data,
	}
}

// Add appends a string to this collector.
func (self Collector) Add(text string) {
	*self.data = append(*self.data, text)
}

// Result provides all accumulated strings.
func (self Collector) Result() []string {
	if self.data == nil {
		return []string{}
	}
	return *self.data
}
