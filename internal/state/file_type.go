package state

type FileType string

const (
	FileTypeRunstate FileType = "runstate"
)

func (self FileType) String() string {
	return string(self)
}
