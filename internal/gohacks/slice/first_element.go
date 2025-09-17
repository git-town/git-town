package slice

import . "github.com/git-town/git-town/v21/pkg/prelude"

func FirstElement[E any, L ~[]E](list L) Option[E] {
	if len(list) > 0 {
		return Some(list[0])
	}
	return None[E]()
}
