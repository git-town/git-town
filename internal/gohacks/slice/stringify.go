package slice

import "fmt"

func Stringify[T fmt.Stringer](elements []T) []string {
	result := make([]string, len(elements))
	for e, element := range elements {
		result[e] = element.String()
	}
	return result
}
