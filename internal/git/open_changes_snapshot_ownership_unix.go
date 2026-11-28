//go:build unix

package git

import (
	"os"
	"syscall"
)

func openChangesCanVerifyFileModes() bool {
	return true
}

func openChangesFileOwnership(info os.FileInfo) (openChangesOwnership, bool) {
	stat, isStat := info.Sys().(*syscall.Stat_t)
	if !isStat {
		return openChangesOwnership{
			groupID:    0,
			hasGroupID: false,
			hasUserID:  false,
			userID:     0,
		}, false
	}
	return openChangesOwnership{
		groupID:    stat.Gid,
		hasGroupID: true,
		hasUserID:  true,
		userID:     stat.Uid,
	}, true
}

func openChangesLchown(path string, userID, groupID uint32, changeUserID, changeGroupID bool) error {
	chownUserID := -1
	if changeUserID {
		chownUserID = int(userID)
	}
	chownGroupID := -1
	if changeGroupID {
		chownGroupID = int(groupID)
	}
	return os.Lchown(path, chownUserID, chownGroupID)
}
