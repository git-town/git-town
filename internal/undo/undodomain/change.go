package undodomain

type Change[T any] struct {
	Before T
	After  T
}
