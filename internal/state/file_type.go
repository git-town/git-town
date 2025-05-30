package state

// FileType defines the types of files and their filenames that can be stored in the persistent state folder on disk.
type FileType string

const (
	FileTypeRunstate FileType = "runstate" // the runstate file
	FileTypeRunlog   FileType = "runlog"   // the runlog file
)

func (self FileType) String() string {
	return string(self)
}
