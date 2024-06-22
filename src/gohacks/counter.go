package gohacks

type Counter int

func NewCounter(value int) Counter {
	return Counter(value)
}

func (self *Counter) Inc() {
	*self += 1
}
