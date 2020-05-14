package test

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// CopyDirectory copies all files in the given src directory into the given dst directory.
// Both the source and the destination directory must exist.
func CopyDirectory(src, dst string) error {
	return filepath.Walk(src, func(srcPath string, fi os.FileInfo, e error) error {
		dstPath := strings.Replace(srcPath, src, dst, 1)
		// handle directory
		if fi.IsDir() {
			err := os.Mkdir(dstPath, fi.Mode())
			if err != nil {
				return fmt.Errorf("cannot create target directory: %w", err)
			}
			return nil
		}
		// handle file
		sourceContent, err := os.Open(srcPath)
		if err != nil {
			return fmt.Errorf("cannot read source file %q: %w", srcPath, err)
		}
		destFile, err := os.OpenFile(dstPath, os.O_CREATE|os.O_WRONLY, fi.Mode())
		if err != nil {
			return fmt.Errorf("cannot create target file %q: %w", srcPath, err)
		}
		_, err = io.Copy(destFile, sourceContent)
		if err != nil {
			return fmt.Errorf("cannot copy %q into %q: %w", srcPath, dstPath, err)
		}
		err = sourceContent.Close()
		if err != nil {
			return fmt.Errorf("cannot close source file %q: %w", srcPath, err)
		}
		err = destFile.Close()
		return err
	})
}

// createFile creates a file with the given filename in the given directory.
func createFile(t *testing.T, dir, filename string) {
	filePath := filepath.Join(dir, filename)
	err := os.MkdirAll(filepath.Dir(filePath), 0744)
	assert.Nil(t, err)
	err = ioutil.WriteFile(filePath, []byte(filename+" content"), 0644)
	assert.Nil(t, err)
}

// createTempDir creates a new empty directory in the system's temp directory and provides the path to it.
func createTempDir(t *testing.T) string {
	dir, err := ioutil.TempDir("", "")
	assert.Nil(t, err, "cannot create TempDir")
	evalDir, err := filepath.EvalSymlinks(dir) // Evaluate symlinks as Mac temp dir is symlinked
	assert.Nil(t, err, "cannot evaluate symlinks of TempDir")
	return evalDir
}
