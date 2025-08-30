package storage

import (
	"io/fs"
	"os"
)

// FS is a wrapper that falls back to `os`.
type FS struct {
	FS fs.FS
}

// Open a file in the provided `fs.FS`. If none provided,
// open via `os.Open`
func (f FS) Open(name string) (fs.File, error) {
	if f.FS == nil {
		return os.Open(name)
	}

	return f.FS.Open(name)
}
