package slice

import . "github.com/git-town/git-town/v16/pkg/prelude"

// returns a copy of given slice where all options are unwrapped and None options are removed
func GetAll[T any](slice []Option[T]) []T {
	result := make([]T, 0, len(slice))
	for _, elementOpt := range slice {
		if element, hasElement := elementOpt.Get(); hasElement {
			result = append(result, element)
		}
	}
	return result
}
