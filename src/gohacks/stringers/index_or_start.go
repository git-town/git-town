package stringers

import "fmt"

func IndexOrStart[S fmt.Stringer](elements []S, needle S) int {
	result := Index(elements, needle)
	if result >= 0 {
		return result
	}
	return 0
}
