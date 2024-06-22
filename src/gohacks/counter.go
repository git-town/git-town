package gohacks

type Counter int

func NewCounter(value int) Counter {
	return Counter(value)
}

// adds 1 to this counter
func (self *Counter) Inc() {
	*self += 1
}
