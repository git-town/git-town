package test

import (
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

// CopyDirectory copies all files in the given src dirctory into the given dst directory.
// Both the source and the destination directory must exist.
func CopyDirectory(src, dst string) error {
	return filepath.Walk(src, func(srcPath string, fi os.FileInfo, err error) error {
		dstPath := strings.Replace(srcPath, src, dst, 1)

		// handle directory
		if fi.IsDir() {
			err := os.Mkdir(dstPath, fi.Mode())
			if err != nil {
				return errors.Wrap(err, "cannot create target directory")
			}
			return nil
		}

		// handle file
		sourceContent, err := os.Open(srcPath)
		defer sourceContent.Close()
		if err != nil {
			return errors.Wrapf(err, "cannot read source file %q", srcPath)
		}
		destFile, err := os.OpenFile(dstPath, os.O_CREATE|os.O_WRONLY, fi.Mode())
		defer destFile.Close()
		if err != nil {
			return errors.Wrapf(err, "Cannot create target file %q", srcPath)
		}
		_, err = io.Copy(destFile, sourceContent)
		if err != nil {
			return errors.Wrapf(err, "cannot copy %q into %q", srcPath, dstPath)
		}
		return nil
	})
}

// createFile creates a file with the given filename in the given directory.
func createFile(t *testing.T, dir, filename string) {
	filePath := path.Join(dir, filename)
	err := os.MkdirAll(path.Dir(filePath), 0744)
	assert.Nil(t, err)
	err = ioutil.WriteFile(filePath, []byte(filename+" content"), 0644)
	assert.Nil(t, err)
}

// createTestDir provides the path to an empty directory in the system's temp directory
func createTempDir(t *testing.T) string {
	dir, err := ioutil.TempDir("", "")
	assert.Nil(t, err, "cannot create TempDir")
	return dir
}
