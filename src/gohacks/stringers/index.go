package stringers

import "fmt"

func Index[S fmt.Stringer](elements []S, needle S) int {
	for e, element := range elements {
		if element.String() == needle.String() {
			return e
		}
	}
	return -1
}
