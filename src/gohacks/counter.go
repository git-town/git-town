package gohacks

type Counter int

// creates a new Counter instance with the given value
func NewCounter(value int) Counter {
	return Counter(value)
}

// adds 1 to this counter
func (self *Counter) Inc() {
	*self += 1
}
