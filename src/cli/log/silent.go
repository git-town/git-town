package log

// The silent logger acts as a stand-in for loggers when no logging is desired.
type Silent struct{}

func (self Silent) Failed(error)                 {}
func (self Silent) Start(string, ...interface{}) {}
func (self Silent) Success()                     {}
