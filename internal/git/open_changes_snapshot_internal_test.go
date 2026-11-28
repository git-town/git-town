package git

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/shoenig/test/must"
)

func TestChmodFileExactWith(t *testing.T) {
	t.Parallel()
	if !openChangesCanVerifyFileModes() {
		t.Skip("this platform does not preserve POSIX permission bits")
	}
	path := filepath.Join(t.TempDir(), "file")
	must.NoError(t, os.WriteFile(path, []byte("content"), 0o600))
	stripSetgid := func(path string, mode os.FileMode) error {
		return os.Chmod(path, mode&^os.ModeSetgid)
	}
	err := chmodFileExactWith(path, 0o644|os.ModeSetgid, stripSetgid)
	must.Error(t, err)
	expected := fmt.Sprintf("cannot apply expected permission bits 2644 to %q: the operating system reports permission bits 0644 after chmod", path)
	must.EqOp(t, expected, err.Error())
}

func TestOpenChangesMetadataID(t *testing.T) {
	t.Parallel()
	id, present, valid := openChangesMetadataID("3000000000")
	must.True(t, valid)
	must.True(t, present)
	must.EqOp(t, uint32(3000000000), id)
	id, present, valid = openChangesMetadataID("-1")
	must.True(t, valid)
	must.False(t, present)
	must.EqOp(t, uint32(0), id)
	_, _, valid = openChangesMetadataID("4294967296")
	must.False(t, valid)
	_, _, valid = openChangesMetadataID("-2")
	must.False(t, valid)
}

func TestRestoreFileExactWithChangedGroup(t *testing.T) {
	t.Parallel()
	if !openChangesCanVerifyFileModes() {
		t.Skip("this platform does not preserve POSIX permission bits")
	}
	tempDir := t.TempDir()
	source := filepath.Join(tempDir, "source")
	destination := filepath.Join(tempDir, "destination")
	must.NoError(t, os.WriteFile(source, []byte("content"), 0o600))
	const recordedGroupID uint32 = 3000000000
	currentOwnership := openChangesOwnership{
		groupID:    20,
		hasGroupID: true,
		hasUserID:  true,
		userID:     10,
	}
	chownCalls := 0
	chown := func(_ string, _, groupID uint32, changeUserID, changeGroupID bool) error {
		chownCalls++
		if changeUserID || !changeGroupID || groupID != recordedGroupID {
			return errors.New("unexpected ownership change")
		}
		currentOwnership.groupID = groupID
		return nil
	}
	fileOwnership := func(os.FileInfo) (openChangesOwnership, bool) {
		return currentOwnership, true
	}
	preserveOwnership := func(path string, metadata openChangesFileMetadata) error {
		return preserveSetIDOwnershipWith(path, metadata, chown, fileOwnership)
	}
	metadata := openChangesFileMetadata{
		mode: 0o644 | os.ModeSetgid,
		ownership: openChangesOwnership{
			groupID:    recordedGroupID,
			hasGroupID: true,
		},
	}
	err := restoreFileExactWith(source, destination, metadata, preserveOwnership)
	must.NoError(t, err)
	must.EqOp(t, 1, chownCalls)
	must.EqOp(t, recordedGroupID, currentOwnership.groupID)
	content, err := os.ReadFile(destination)
	must.NoError(t, err)
	must.EqOp(t, "content", string(content))
	info, err := os.Lstat(destination)
	must.NoError(t, err)
	must.EqOp(t, os.FileMode(0o644), info.Mode().Perm())
	must.True(t, info.Mode()&os.ModeSetgid != 0)
}

func TestRestoreFileExactWithChangedOwner(t *testing.T) {
	t.Parallel()
	if !openChangesCanVerifyFileModes() {
		t.Skip("this platform does not preserve POSIX permission bits")
	}
	tempDir := t.TempDir()
	source := filepath.Join(tempDir, "source")
	destination := filepath.Join(tempDir, "destination")
	must.NoError(t, os.WriteFile(source, []byte("content"), 0o600))
	const recordedUserID uint32 = 3000000000
	currentOwnership := openChangesOwnership{
		groupID:    20,
		hasGroupID: true,
		hasUserID:  true,
		userID:     10,
	}
	chownCalls := 0
	chown := func(path string, userID, _ uint32, changeUserID, changeGroupID bool) error {
		chownCalls++
		if !changeUserID || changeGroupID || userID != recordedUserID {
			return errors.New("unexpected ownership change")
		}
		info, err := os.Lstat(path)
		if err != nil {
			return err
		}
		if info.Mode()&os.ModeSetuid != 0 {
			return errors.New("owner changed after chmod")
		}
		currentOwnership.userID = userID
		return nil
	}
	fileOwnership := func(os.FileInfo) (openChangesOwnership, bool) {
		return currentOwnership, true
	}
	preserveOwnership := func(path string, metadata openChangesFileMetadata) error {
		return preserveSetIDOwnershipWith(path, metadata, chown, fileOwnership)
	}
	metadata := openChangesFileMetadata{
		mode: 0o644 | os.ModeSetuid,
		ownership: openChangesOwnership{
			hasUserID: true,
			userID:    recordedUserID,
		},
	}
	err := restoreFileExactWith(source, destination, metadata, preserveOwnership)
	must.NoError(t, err)
	must.EqOp(t, 1, chownCalls)
	must.EqOp(t, recordedUserID, currentOwnership.userID)
	content, err := os.ReadFile(destination)
	must.NoError(t, err)
	must.EqOp(t, "content", string(content))
	info, err := os.Lstat(destination)
	must.NoError(t, err)
	must.EqOp(t, os.FileMode(0o644), info.Mode().Perm())
	must.True(t, info.Mode()&os.ModeSetuid != 0)
}
