//go:build !unix

package git

import (
	"errors"
	"os"
)

func openChangesCanVerifyFileModes() bool {
	return false
}

func openChangesFileOwnership(os.FileInfo) (openChangesOwnership, bool) {
	return openChangesOwnership{
		groupID:    0,
		hasGroupID: false,
		hasUserID:  false,
		userID:     0,
	}, false
}

func openChangesLchown(string, uint32, uint32, bool, bool) error {
	return errors.New("the operating system does not expose file ownership")
}
