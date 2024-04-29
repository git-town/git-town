package prelude

// SafePointer is a pointer that is guaranteed to not be nil.
// This helps remove the possibility for nil-pointer-panics at runtime.
// It also makes the code cleaner by removing unnecessary defensive nil checks
// "just to be careful" for pointers that are known to never be nil.
// You must create new instances via NewSafePointer.
// The Zero value is an intentionally uninitialized SafePointer that panics when used.
type SafePointer[T any] struct {
	value       *T
	initialized bool
}

func NewSafePointer[T any](value *T) SafePointer[T] {
	if value == nil {
		panic("cannot initialize SafePointer with nil")
	}
	return SafePointer[T]{value, true}
}

func (self SafePointer[T]) Get() *T {
	if !self.initialized {
		panic("cannot use an uninitialized SafePointer")
	}
	return self.value
}
