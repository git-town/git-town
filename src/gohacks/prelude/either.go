package prelude

// Either is a type that contains either of the two types, but never both.
type Either[LEFT any, RIGHT any] struct {
	left  *LEFT
	right *RIGHT
}

func Left[LEFT any, RIGHT any](left LEFT) Either[LEFT, RIGHT] {
	return Either[LEFT, RIGHT]{
		left:  &left,
		right: nil,
	}
}

func Right[LEFT any, RIGHT any](right RIGHT) Either[LEFT, RIGHT] {
	return Either[LEFT, RIGHT]{
		left:  nil,
		right: &right,
	}
}

func (self Either[LEFT, RIGHT]) Get() (left LEFT, hasLeft bool, right RIGHT, hasRight bool) {
	if self.left != nil {
		left = *self.left
		hasLeft = true
		return
	}
	right = *self.right
	hasRight = true
	return
}
