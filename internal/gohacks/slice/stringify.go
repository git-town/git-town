package slice

import "fmt"

// provides the string version of the given stringers
func Stringify[T fmt.Stringer](elements []T) []string {
	result := make([]string, len(elements))
	for e, element := range elements {
		result[e] = element.String()
	}
	return result
}
