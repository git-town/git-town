package gohacks

import (
	"fmt"
	"os"
)

// ReadTextFilePanic provides the content of the given file.
// Panics if there is an error, so only use for debugging.
func ReadTextFilePanic(path string) string {
	content, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Sprintf("cannot read file %q: %s", path, err))
	}
	return string(content)
}
