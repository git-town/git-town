package state

// FileType defines the types of files and their filenames that can be stored in the persistent state folder on disk.
type FileType string

const (
	FileTypeRunstate FileType = "runstate"
)

func (self FileType) String() string {
	return string(self)
}
