package gohacks

// Counter is a Statistics implementation that counts how many commands were run.
// This struct is always mutable, even if you provide or call it by value.
type Counter struct {
	count *int
}

func NewCounter() Counter {
	var count int
	return Counter{
		count: &count,
	}
}

func (self *Counter) Count() int {
	return *self.count
}

func (self *Counter) Register() {
	*self.count++
}
