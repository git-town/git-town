package prelude

// Either is a type that contains either of the two types, but never both.
type Either[MAIN any, OTHER any] struct {
	main  *MAIN
	other *OTHER
}

func Main[MAIN any, OTHER any](main MAIN) Either[MAIN, OTHER] {
	return Either[MAIN, OTHER]{
		main:  &main,
		other: nil,
	}
}

func Other[MAIN any, OTHER any](other OTHER) Either[MAIN, OTHER] {
	return Either[MAIN, OTHER]{
		main:  nil,
		other: &other,
	}
}

func (self Either[MAIN, OTHER]) Get() (main MAIN, hasMain bool, other OTHER, hasOther bool) {
	if self.main != nil {
		main = *self.main
		hasMain = true
		return
	}
	other = *self.other
	hasOther = true
	return
}
