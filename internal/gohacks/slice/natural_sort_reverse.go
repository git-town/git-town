package slice

import (
	"fmt"
	"slices"
)

// sorts the given elements in reverse natural sort order
func NaturalSortReverse[T fmt.Stringer](list []T) {
	if len(list) < 2 {
		return
	}
	NaturalSort(list)
	slices.Reverse(list)
}
