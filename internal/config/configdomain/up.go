package configdomain

// Up indicates that "git town commit" should commit into the parent branch
type Up uint

func (self Up) Value() uint {
	return uint(self)
}
