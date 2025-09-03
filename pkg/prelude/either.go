package prelude

// Either is a type that contains either of the two types, but never both.
type Either[LEFT any, RIGHT any] struct {
	left  *LEFT
	right *RIGHT
}

// Left creates an Either containing the given left value
func Left[LEFT any, RIGHT any](left LEFT) Either[LEFT, RIGHT] {
	return Either[LEFT, RIGHT]{
		left:  &left,
		right: nil,
	}
}

// Right creates an Either containing the given right value
func Right[LEFT any, RIGHT any](right RIGHT) Either[LEFT, RIGHT] {
	return Either[LEFT, RIGHT]{
		left:  nil,
		right: &right,
	}
}

// Get returns the contained value and indicates which side it is on
func (self Either[LEFT, RIGHT]) Get() (left LEFT, hasLeft bool, right RIGHT, hasRight bool) { //nolint:ireturn
	if self.left != nil {
		left = *self.left
		hasLeft = true
		return
	}
	right = *self.right
	hasRight = true
	return
}
