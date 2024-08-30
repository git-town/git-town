package slice

import . "github.com/git-town/git-town/v16/pkg/prelude"

func Index[E comparable, C ~[]E](haystack C, needle E) Option[int] {
	for e, element := range haystack {
		if element == needle {
			return Some(e)
		}
	}
	return None[int]()
}
