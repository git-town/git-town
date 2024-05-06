package stringslice

// Collector accumulates string instances one at a time.
// The zero value is an empty collection.
type Collector struct {
	data *[]string
}

func NewCollector() Collector {
	data := []string{}
	return Collector{
		data: &data,
	}
}

// Add appends a string instance to this collector.
func (self Collector) Add(text string) {
	*self.data = append(*self.data, text)
}

// Result provides all accumulated string instances.
func (self Collector) Result() []string {
	if self.data == nil {
		return []string{}
	}
	return *self.data
}
