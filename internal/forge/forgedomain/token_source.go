package forgedomain

type TokenSource string

const (
	TokenSourceManual TokenSource = "manual"
	TokenSourceScript TokenSource = "script"
)

func (self TokenSource) String() string {
	return string(self)
}
