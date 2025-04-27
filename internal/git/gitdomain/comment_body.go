package gitdomain

type CommentBody string

func (self CommentBody) String() string {
	return string(self)
}
