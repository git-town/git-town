package configdomain

// Down indicates that "git town commit" should commit into the parent branch
type Down int

func (self Down) Value() int {
	return int(self)
}
