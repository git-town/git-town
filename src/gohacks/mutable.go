package gohacks

type Mutable[T any] struct {
	instance *T
}

func (self Mutable[T]) Get() *T {
	if self.instance == nil {
		self.instance = new(T)
	}
	return self.instance
}
