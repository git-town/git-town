package filesystem

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/git-town/git-town/v14/test/asserts"
)

// CopyDirectory copies all files in the given src directory into the given dst directory.
// Both the source and the destination directory must exist.
func CopyDirectory(src, dst string) {
	asserts.NoError(filepath.Walk(src, func(srcPath string, fileInfo os.FileInfo, _ error) error {
		dstPath := strings.Replace(srcPath, src, dst, 1)
		if fileInfo.IsDir() {
			err := os.Mkdir(dstPath, fileInfo.Mode())
			if err != nil {
				return fmt.Errorf("cannot create target directory: %w", err)
			}
			return nil
		}
		sourceFile, err := os.Open(srcPath)
		if err != nil {
			return fmt.Errorf("cannot open source file %q: %w", srcPath, err)
		}
		destFile, err := os.OpenFile(dstPath, os.O_CREATE|os.O_WRONLY, fileInfo.Mode())
		if err != nil {
			return fmt.Errorf("cannot create target file %q: %w", srcPath, err)
		}
		_, err = io.Copy(destFile, sourceFile)
		if err != nil {
			return fmt.Errorf("cannot copy %q into %q: %w", srcPath, dstPath, err)
		}
		err = sourceFile.Close()
		if err != nil {
			return fmt.Errorf("cannot close source file %q: %w", srcPath, err)
		}
		err = destFile.Close()
		return err
	}))
}
