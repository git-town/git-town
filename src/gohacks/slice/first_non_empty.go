package slice

type emptyable interface {
	IsEmpty() bool
}

// FirstNonEmpty provides the first of its arguments that isn't empty.
func FirstNonEmpty[T emptyable](first T, others ...T) T {
	if !first.IsEmpty() {
		return first
	}
	for _, other := range others {
		if !other.IsEmpty() {
			return other
		}
	}
	panic("only empty elements")
}
