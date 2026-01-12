package configdomain

// Down indicates that "git town commit" should commit into the parent branch
type Down uint

func (self Down) Value() uint {
	return uint(self)
}
