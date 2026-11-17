package undodomain

type Change[T any] struct {
	After  T
	Before T
}
