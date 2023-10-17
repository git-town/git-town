package log

type Silent struct{}

func (self Silent) Start(string, ...interface{}) {}
func (self Silent) Success()                     {}
func (self Silent) Failed(error)                 {}
