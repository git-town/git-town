package gohacks

// Counter is a Statistics implementation that counts how many commands were run.
type Counter struct {
	count *int
}

func (self *Counter) Count() int {
	return *self.count
}

func (self *Counter) Register() {
	*self.count++
}
