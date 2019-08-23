package infra

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

// CopyDirectory copies all files in the given src dirctory into the given dst directory.
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
			return errors.Wrapf(err, "cannot read source file '%s'", srcPath)
		}
		destFile, err := os.OpenFile(dstPath, os.O_CREATE|os.O_WRONLY, fi.Mode())
		defer destFile.Close()
		if err != nil {
			return errors.Wrapf(err, "Cannot create target file '%s'", srcPath)
		}
		_, err = io.Copy(destFile, sourceContent)
		if err != nil {
			return errors.Wrapf(err, "cannot copy '%s' into '%s'", srcPath, dstPath)
		}
		return nil
	})
}
