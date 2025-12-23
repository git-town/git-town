package stringslice

import "fmt"

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
// TODO: add an AddF method that receives a format string and a variadic number of arguments.
func (self Collector) Add(text string) {
	*self.data = append(*self.data, text)
}

func (self Collector) Addf(format string, args ...any) {
	self.Add(fmt.Sprintf(format, args...))
}

// Result provides all accumulated strings.
func (self Collector) Result() []string {
	if self.data == nil {
		return []string{}
	}
	return *self.data
}
