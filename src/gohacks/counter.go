package gohacks

// Counter is a special type used for counting things.
// The zero value is a valid empty counter.
type Counter int

// adds 1 to this counter
func (self *Counter) Inc() {
	*self += 1
}
