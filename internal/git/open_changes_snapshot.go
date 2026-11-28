package git

import (
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/git-town/git-town/v24/internal/git/gitdomain"
	"github.com/git-town/git-town/v24/internal/gohacks/stringss"
	"github.com/git-town/git-town/v24/internal/messages"
	"github.com/git-town/git-town/v24/internal/subshell/subshelldomain"
	. "github.com/git-town/git-town/v24/pkg/prelude"
)

const (
	openChangesBackupRoot          = "git-town/open-changes"
	openChangesSnapshotRefPrefix   = "refs/git-town/open-changes/"
	openChangesRawFilesCommitTitle = "raw files in Git Town open changes snapshot"
	openChangesRawFileModesHeader  = "git-town-open-changes-format: modes-v2"
	emptyTreeSHA                   = "4b825dc642cb6eb9a060e54bf8d69288fbee4904"
)

// CreateOpenChangesSnapshot creates a reachable stash-shaped commit that captures the exact index and worktree state,
// including raw file bytes that Git filters or end-of-line conversion cannot represent, without modifying either one.
func (self *Commands) CreateOpenChangesSnapshot(backend subshelldomain.RunnerQuerier) (Option[gitdomain.SHA], error) {
	trackedSnapshotText, err := backend.QueryTrim("git", "-c", "core.safecrlf=false", "stash", "create")
	if err != nil {
		return None[gitdomain.SHA](), fmt.Errorf(messages.OpenChangesSnapshotCreateProblem, err)
	}
	untrackedFiles, err := listUntrackedFiles(backend)
	if err != nil {
		return None[gitdomain.SHA](), err
	}
	rawFiles := append([]string{}, untrackedFiles...)
	trackedSnapshot := None[gitdomain.SHA]()
	if len(trackedSnapshotText) > 0 {
		parsedTrackedSnapshot, err := gitdomain.NewSHA(trackedSnapshotText)
		if err != nil {
			return None[gitdomain.SHA](), fmt.Errorf(messages.OpenChangesSnapshotParseProblem, trackedSnapshotText, err)
		}
		trackedSnapshot = Some(parsedTrackedSnapshot)
	}
	if parsedTrackedSnapshot, hasTrackedSnapshot := trackedSnapshot.Get(); hasTrackedSnapshot {
		trackedOpenFiles, err := listTrackedOpenChangeFiles(backend)
		if err != nil {
			return None[gitdomain.SHA](), err
		}
		if len(trackedOpenFiles) > 0 {
			rootDir, hasRootDir := self.RootDirectory(backend).Get()
			if !hasRootDir {
				return None[gitdomain.SHA](), errors.New(messages.OpenChangesSnapshotCreateNoRoot)
			}
			trackedRawFiles, err := self.trackedFilesNeedingRawBackup(backend, rootDir, parsedTrackedSnapshot, trackedOpenFiles)
			if err != nil {
				return None[gitdomain.SHA](), err
			}
			rawFiles = mergeOpenChangesFileLists(rawFiles, trackedRawFiles)
		}
	}
	var backup *openChangesBackup
	backupComplete := false
	defer func() {
		if backup != nil && !backupComplete {
			backup.discard()
		}
	}()
	rawFilesCommit := None[gitdomain.SHA]()
	if len(rawFiles) > 0 {
		rootDir, hasRootDir := self.RootDirectory(backend).Get()
		if !hasRootDir {
			return None[gitdomain.SHA](), errors.New(messages.OpenChangesSnapshotCreateNoRoot)
		}
		gitDir, err := self.gitDirectory(backend)
		if err != nil {
			return None[gitdomain.SHA](), fmt.Errorf(messages.OpenChangesSnapshotCreateProblem, err)
		}
		createdBackup, err := createOpenChangesBackup(rootDir, gitDir, rawFiles)
		if err != nil {
			return None[gitdomain.SHA](), err
		}
		backup = &createdBackup
		rawFilesCommit, err = self.createRawFilesCommit(backend, createdBackup)
		if err != nil {
			return None[gitdomain.SHA](), err
		}
	}
	var snapshot gitdomain.SHA
	if parsedTrackedSnapshot, hasTrackedSnapshot := trackedSnapshot.Get(); !hasTrackedSnapshot {
		if rawFilesCommit.IsNone() {
			return None[gitdomain.SHA](), nil
		}
		snapshot, err = self.createRawFilesOnlySnapshot(backend, rawFilesCommit.GetOrPanic())
	} else {
		if rawFilesCommit.IsNone() {
			snapshot = parsedTrackedSnapshot
		} else {
			snapshot, err = self.addRawFilesToSnapshot(backend, parsedTrackedSnapshot, rawFilesCommit.GetOrPanic())
		}
	}
	if err != nil {
		return None[gitdomain.SHA](), err
	}
	if backup != nil {
		gitDir, err := self.gitDirectory(backend)
		if err != nil {
			return None[gitdomain.SHA](), fmt.Errorf(messages.OpenChangesSnapshotCreateProblem, err)
		}
		if err := backup.finalize(gitDir, snapshot); err != nil {
			return None[gitdomain.SHA](), fmt.Errorf(messages.OpenChangesSnapshotCreateProblem, err)
		}
	}
	if err := backend.Run("git", "update-ref", OpenChangesSnapshotRef(snapshot), snapshot.String()); err != nil {
		return None[gitdomain.SHA](), fmt.Errorf(messages.OpenChangesSnapshotCreateProblem, err)
	}
	backupComplete = true
	return Some(snapshot), nil
}

// OpenChangesSnapshotRef provides the ref that keeps the given open changes snapshot reachable.
func OpenChangesSnapshotRef(snapshot gitdomain.SHA) string {
	return openChangesSnapshotRefPrefix + snapshot.String()
}

// RestoreOpenChangesSnapshot applies the given snapshot, restores the raw file bytes atomically,
// and removes the backup data once it is no longer needed.
func (self *Commands) RestoreOpenChangesSnapshot(backend subshelldomain.RunnerQuerier, frontend subshelldomain.Runner, snapshot gitdomain.SHA) error {
	if err := self.ApplyStash(frontend, snapshot); err != nil {
		return fmt.Errorf(messages.OpenChangesSnapshotRestoreProblem, err)
	}
	rawFilesCommit, err := openChangesSnapshotRawCommit(backend, snapshot)
	if err != nil {
		return fmt.Errorf(messages.OpenChangesSnapshotRestoreProblem, err)
	}
	if rawFilesCommit, hasRawFilesCommit := rawFilesCommit.Get(); hasRawFilesCommit {
		gitDir, err := self.gitDirectory(backend)
		if err != nil {
			return fmt.Errorf(messages.OpenChangesSnapshotRestoreProblem, err)
		}
		rootDir, hasRootDir := self.RootDirectory(backend).Get()
		if !hasRootDir {
			return errors.New(messages.OpenChangesSnapshotRestoreNoRoot)
		}
		entries, err := openChangesTreeEntries(backend, rawFilesCommit)
		if err != nil {
			return fmt.Errorf(messages.OpenChangesSnapshotRestoreProblem, err)
		}
		rawFilesMetadata, validRawFilesMetadata, err := openChangesRawFilesMetadata(backend, rawFilesCommit)
		if err != nil {
			return fmt.Errorf(messages.OpenChangesSnapshotRestoreProblem, err)
		}
		if !validRawFilesMetadata || !openChangesMetadataMatchEntries(entries, rawFilesMetadata) {
			return errors.New(messages.OpenChangesSnapshotInvalid)
		}
		backupDir := openChangesBackupDir(gitDir, snapshot)
		for _, file := range sortedOpenChangesTreeFiles(entries) {
			source := filepath.Join(backupDir, "files", filepath.FromSlash(file))
			destination, err := worktreeFilePath(rootDir, file)
			if err != nil {
				return err
			}
			if err := restoreFileExact(source, destination, rawFilesMetadata[file]); err != nil {
				return fmt.Errorf(messages.OpenChangesSnapshotRestoreProblem, err)
			}
		}
		if err := os.RemoveAll(backupDir); err != nil {
			return fmt.Errorf(messages.OpenChangesSnapshotRestoreProblem, err)
		}
	}
	if err := frontend.Run("git", "update-ref", "-d", OpenChangesSnapshotRef(snapshot)); err != nil {
		return fmt.Errorf(messages.OpenChangesSnapshotRestoreProblem, err)
	}
	return nil
}

// ValidateOpenChangesSnapshot verifies that the snapshot and every file in its raw backup are still complete and unchanged.
func (self *Commands) ValidateOpenChangesSnapshot(backend subshelldomain.RunnerQuerier, snapshot gitdomain.SHA) error {
	actualSHA, err := backend.QueryTrim("git", "rev-parse", "--verify", "-q", OpenChangesSnapshotRef(snapshot))
	if err != nil || actualSHA.String() != snapshot.String() {
		return errors.New(messages.OpenChangesSnapshotInvalid)
	}
	rawFilesCommit, err := openChangesSnapshotRawCommit(backend, snapshot)
	if err != nil {
		return fmt.Errorf(messages.OpenChangesSnapshotValidateProblem, err)
	}
	if rawFilesCommit.IsNone() {
		return nil
	}
	gitDir, err := self.gitDirectory(backend)
	if err != nil {
		return fmt.Errorf(messages.OpenChangesSnapshotValidateProblem, err)
	}
	backupDir := openChangesBackupDir(gitDir, snapshot)
	if info, err := os.Stat(backupDir); err != nil || !info.IsDir() {
		return errors.New(messages.OpenChangesSnapshotInvalid)
	}
	entries, err := openChangesTreeEntries(backend, rawFilesCommit.GetOrPanic())
	if err != nil {
		return fmt.Errorf(messages.OpenChangesSnapshotValidateProblem, err)
	}
	if len(entries) == 0 {
		return errors.New(messages.OpenChangesSnapshotInvalid)
	}
	rawFilesMetadata, validRawFilesMetadata, err := openChangesRawFilesMetadata(backend, rawFilesCommit.GetOrPanic())
	if err != nil {
		return fmt.Errorf(messages.OpenChangesSnapshotValidateProblem, err)
	}
	if !validRawFilesMetadata || !openChangesMetadataMatchEntries(entries, rawFilesMetadata) {
		return errors.New(messages.OpenChangesSnapshotInvalid)
	}
	for _, file := range sortedOpenChangesTreeFiles(entries) {
		entry := entries[file]
		if entry.objectType != "blob" {
			return errors.New(messages.OpenChangesSnapshotInvalid)
		}
		backupPath := filepath.Join(backupDir, "files", filepath.FromSlash(file))
		info, err := os.Lstat(backupPath)
		if err != nil || !openChangesBackupTypeMatches(info, entry.permission) {
			return errors.New(messages.OpenChangesSnapshotInvalid)
		}
		metadata := rawFilesMetadata[file]
		if openChangesRegularPermission(entry.permission) && openChangesPermissionBits(info.Mode()) != openChangesPermissionBits(metadata.mode) {
			return errors.New(messages.OpenChangesSnapshotInvalid)
		}
		if openChangesRegularPermission(entry.permission) && !openChangesBackupOwnershipMatches(info, metadata) {
			return errors.New(messages.OpenChangesSnapshotInvalid)
		}
		actualSHA, err := hashBackupFileExact(backend, backupPath, info)
		if err != nil {
			return fmt.Errorf(messages.OpenChangesSnapshotValidateProblem, err)
		}
		if actualSHA != entry.sha {
			return errors.New(messages.OpenChangesSnapshotInvalid)
		}
	}
	return nil
}

func (self *Commands) addRawFilesToSnapshot(backend subshelldomain.RunnerQuerier, trackedSnapshot, rawFilesCommit gitdomain.SHA) (gitdomain.SHA, error) {
	tree, err := querySHA(backend, "rev-parse", trackedSnapshot.String()+"^{tree}")
	if err != nil {
		return "", fmt.Errorf(messages.OpenChangesSnapshotCreateProblem, err)
	}
	head, err := querySHA(backend, "rev-parse", "HEAD")
	if err != nil {
		return "", fmt.Errorf(messages.OpenChangesSnapshotCreateProblem, err)
	}
	indexCommit, err := querySHA(backend, "rev-parse", trackedSnapshot.String()+"^2")
	if err != nil {
		return "", fmt.Errorf(messages.OpenChangesSnapshotCreateProblem, err)
	}
	emptyCommit, err := createEmptyOpenChangesCommit(backend)
	if err != nil {
		return "", err
	}
	snapshot, err := querySHA(backend, "commit-tree", tree.String(), "-p", head.String(), "-p", indexCommit.String(), "-p", emptyCommit.String(), "-p", rawFilesCommit.String(), "-m", "Git Town open changes snapshot")
	if err != nil {
		return "", fmt.Errorf(messages.OpenChangesSnapshotCreateProblem, err)
	}
	return snapshot, nil
}

func chmodFileExact(path string, expectedMode os.FileMode) error {
	return chmodFileExactWith(path, expectedMode, os.Chmod)
}

func chmodFileExactWith(path string, expectedMode os.FileMode, chmod func(string, os.FileMode) error) error {
	if err := chmod(path, expectedMode); err != nil {
		return err
	}
	if !openChangesCanVerifyFileModes() {
		return nil
	}
	info, err := os.Lstat(path)
	if err != nil {
		return err
	}
	expectedBits := openChangesPermissionBits(expectedMode)
	actualBits := openChangesPermissionBits(info.Mode())
	if actualBits != expectedBits {
		return fmt.Errorf(messages.OpenChangesSnapshotModeMismatch, expectedBits, path, actualBits)
	}
	return nil
}

func copyFileContents(source string, output *os.File) error {
	input, err := os.Open(source)
	if err != nil {
		_ = output.Close()
		return err
	}
	defer input.Close()
	_, copyErr := io.Copy(output, input)
	closeErr := output.Close()
	if copyErr != nil {
		return copyErr
	}
	return closeErr
}

func (self *Commands) createRawFilesCommit(backend subshelldomain.RunnerQuerier, backup openChangesBackup) (Option[gitdomain.SHA], error) {
	indexFile, err := os.CreateTemp(backup.dir, "index-*")
	if err != nil {
		return None[gitdomain.SHA](), fmt.Errorf(messages.OpenChangesSnapshotCreateProblem, err)
	}
	indexPath := indexFile.Name()
	if err := indexFile.Close(); err != nil {
		return None[gitdomain.SHA](), fmt.Errorf(messages.OpenChangesSnapshotCreateProblem, err)
	}
	if err := os.Remove(indexPath); err != nil {
		return None[gitdomain.SHA](), fmt.Errorf(messages.OpenChangesSnapshotCreateProblem, err)
	}
	defer os.Remove(indexPath)
	env := []string{"GIT_INDEX_FILE=" + indexPath}
	modeLines := make([]string, 0, len(backup.files))
	for _, file := range backup.files {
		backupPath := backup.filePath(file)
		info, err := os.Lstat(backupPath)
		if err != nil {
			return None[gitdomain.SHA](), fmt.Errorf(messages.OpenChangesSnapshotCreateProblem, err)
		}
		hashPath := backupPath
		permission := "100644"
		switch {
		case info.Mode()&os.ModeSymlink != 0:
			permission = "120000"
			target, err := os.Readlink(backupPath)
			if err != nil {
				return None[gitdomain.SHA](), fmt.Errorf(messages.OpenChangesSnapshotCreateProblem, err)
			}
			hashPath = filepath.Join(backup.dir, "symlink-content")
			if err := os.WriteFile(hashPath, []byte(target), 0o600); err != nil {
				return None[gitdomain.SHA](), fmt.Errorf(messages.OpenChangesSnapshotCreateProblem, err)
			}
		case info.Mode().Perm()&0o111 != 0:
			permission = "100755"
		}
		if info.Mode().IsRegular() {
			metadata, err := newOpenChangesFileMetadata(backupPath, info)
			if err != nil {
				return None[gitdomain.SHA](), fmt.Errorf(messages.OpenChangesSnapshotCreateProblem, err)
			}
			modeLines = append(modeLines, fmt.Sprintf("mode %04o %s %s %s", openChangesPermissionBits(metadata.mode), openChangesMetadataIDText(metadata.ownership.userID, metadata.ownership.hasUserID), openChangesMetadataIDText(metadata.ownership.groupID, metadata.ownership.hasGroupID), hex.EncodeToString([]byte(file))))
		}
		blob, err := querySHA(backend, "hash-object", "-w", "--literally", "--no-filters", "--", hashPath)
		if info.Mode()&os.ModeSymlink != 0 {
			_ = os.Remove(hashPath)
		}
		if err != nil {
			return None[gitdomain.SHA](), fmt.Errorf(messages.OpenChangesSnapshotCreateProblem, err)
		}
		if err := backend.RunWithEnv(env, "git", "update-index", "--add", "--cacheinfo", permission+","+blob.String()+","+file); err != nil {
			return None[gitdomain.SHA](), fmt.Errorf(messages.OpenChangesSnapshotCreateProblem, err)
		}
	}
	treeOutput, err := backend.QueryWithEnv(env, "git", "write-tree")
	if err != nil {
		return None[gitdomain.SHA](), fmt.Errorf(messages.OpenChangesSnapshotCreateProblem, err)
	}
	tree, err := gitdomain.NewSHA(stringss.Trim(treeOutput))
	if err != nil {
		return None[gitdomain.SHA](), fmt.Errorf(messages.OpenChangesSnapshotParseProblem, treeOutput, err)
	}
	messageFile, err := os.CreateTemp(backup.dir, "message-*")
	if err != nil {
		return None[gitdomain.SHA](), fmt.Errorf(messages.OpenChangesSnapshotCreateProblem, err)
	}
	messagePath := messageFile.Name()
	defer os.Remove(messagePath)
	if _, err := messageFile.WriteString(openChangesRawFilesCommitMessage(modeLines)); err != nil {
		_ = messageFile.Close()
		return None[gitdomain.SHA](), fmt.Errorf(messages.OpenChangesSnapshotCreateProblem, err)
	}
	if err := messageFile.Close(); err != nil {
		return None[gitdomain.SHA](), fmt.Errorf(messages.OpenChangesSnapshotCreateProblem, err)
	}
	commit, err := querySHA(backend, "commit-tree", tree.String(), "-F", messagePath)
	if err != nil {
		return None[gitdomain.SHA](), fmt.Errorf(messages.OpenChangesSnapshotCreateProblem, err)
	}
	return Some(commit), nil
}

func (self *Commands) createRawFilesOnlySnapshot(backend subshelldomain.RunnerQuerier, rawFilesCommit gitdomain.SHA) (gitdomain.SHA, error) {
	head, err := querySHA(backend, "rev-parse", "HEAD")
	if err != nil {
		return "", fmt.Errorf(messages.OpenChangesSnapshotCreateProblem, err)
	}
	headTree, err := querySHA(backend, "rev-parse", "HEAD^{tree}")
	if err != nil {
		return "", fmt.Errorf(messages.OpenChangesSnapshotCreateProblem, err)
	}
	indexCommit, err := querySHA(backend, "commit-tree", headTree.String(), "-p", head.String(), "-m", "index in Git Town open changes snapshot")
	if err != nil {
		return "", fmt.Errorf(messages.OpenChangesSnapshotCreateProblem, err)
	}
	emptyCommit, err := createEmptyOpenChangesCommit(backend)
	if err != nil {
		return "", err
	}
	snapshot, err := querySHA(backend, "commit-tree", headTree.String(), "-p", head.String(), "-p", indexCommit.String(), "-p", emptyCommit.String(), "-p", rawFilesCommit.String(), "-m", "Git Town open changes snapshot")
	if err != nil {
		return "", fmt.Errorf(messages.OpenChangesSnapshotCreateProblem, err)
	}
	return snapshot, nil
}

func (self *Commands) trackedFilesNeedingRawBackup(backend subshelldomain.RunnerQuerier, rootDir gitdomain.RepoRootDir, trackedSnapshot gitdomain.SHA, files []string) ([]string, error) {
	entries, err := openChangesTreeEntries(backend, trackedSnapshot)
	if err != nil {
		return []string{}, fmt.Errorf(messages.OpenChangesSnapshotCreateProblem, err)
	}
	result := make([]string, 0, len(files))
	for _, file := range files {
		entry, hasEntry := entries[file]
		if !hasEntry || entry.objectType != "blob" {
			continue
		}
		source, err := worktreeFilePath(rootDir, file)
		if err != nil {
			return []string{}, err
		}
		info, err := os.Lstat(source)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return []string{}, fmt.Errorf(messages.OpenChangesSnapshotCreateProblem, err)
		}
		if info.Mode()&os.ModeSymlink != 0 || !info.Mode().IsRegular() {
			continue
		}
		result = append(result, file)
	}
	return result, nil
}

type openChangesBackup struct {
	dir   string
	files []string
}

func (self *openChangesBackup) discard() {
	_ = os.RemoveAll(self.dir)
}

func (self *openChangesBackup) filePath(file string) string {
	return filepath.Join(self.dir, "files", filepath.FromSlash(file))
}

func (self *openChangesBackup) finalize(gitDir stringss.Trimmed, snapshot gitdomain.SHA) error {
	target := openChangesBackupDir(gitDir, snapshot)
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return err
	}
	if _, err := os.Stat(target); err == nil {
		self.discard()
		self.dir = target
		return nil
	}
	if err := os.Rename(self.dir, target); err != nil {
		return err
	}
	self.dir = target
	return nil
}

type openChangesFileMetadata struct {
	mode      os.FileMode
	ownership openChangesOwnership
}

type openChangesOwnership struct {
	groupID    uint32
	hasGroupID bool
	hasUserID  bool
	userID     uint32
}

type openChangesTreeEntry struct {
	objectType string
	permission string
	sha        gitdomain.SHA
}

func copyFileExact(source, destination string) error {
	info, err := os.Lstat(source)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(destination), 0o755); err != nil {
		return err
	}
	if info.Mode()&os.ModeSymlink != 0 {
		target, err := os.Readlink(source)
		if err != nil {
			return err
		}
		return os.Symlink(target, destination)
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf(messages.OpenChangesSnapshotUnsupportedFile, source)
	}
	metadata, err := newOpenChangesFileMetadata(source, info)
	if err != nil {
		return err
	}
	output, err := os.OpenFile(destination, os.O_CREATE|os.O_EXCL|os.O_WRONLY, metadata.mode)
	if err != nil {
		return err
	}
	if err := copyFileContents(source, output); err != nil {
		_ = os.Remove(destination)
		return err
	}
	if err := preserveSetIDOwnership(destination, metadata); err != nil {
		_ = os.Remove(destination)
		return err
	}
	if err := chmodFileExact(destination, metadata.mode); err != nil {
		_ = os.Remove(destination)
		return err
	}
	return nil
}

func createEmptyOpenChangesCommit(backend subshelldomain.RunnerQuerier) (gitdomain.SHA, error) {
	commit, err := querySHA(backend, "commit-tree", emptyTreeSHA, "-m", "empty untracked files in Git Town open changes snapshot")
	if err != nil {
		return "", fmt.Errorf(messages.OpenChangesSnapshotCreateProblem, err)
	}
	return commit, nil
}

func createOpenChangesBackup(rootDir gitdomain.RepoRootDir, gitDir stringss.Trimmed, files []string) (openChangesBackup, error) {
	dir, err := os.MkdirTemp(gitDir.String(), "open-changes-*")
	if err != nil {
		return openChangesBackup{}, fmt.Errorf(messages.OpenChangesSnapshotCreateProblem, err)
	}
	result := openChangesBackup{dir: dir, files: files}
	for _, file := range files {
		source, err := worktreeFilePath(rootDir, file)
		if err != nil {
			result.discard()
			return openChangesBackup{}, err
		}
		if err := copyFileExact(source, result.filePath(file)); err != nil {
			result.discard()
			return openChangesBackup{}, fmt.Errorf(messages.OpenChangesSnapshotCreateProblem, err)
		}
	}
	return result, nil
}

func hashBackupFileExact(backend subshelldomain.RunnerQuerier, backupPath string, info os.FileInfo) (gitdomain.SHA, error) {
	hashPath := backupPath
	if info.Mode()&os.ModeSymlink != 0 {
		target, err := os.Readlink(backupPath)
		if err != nil {
			return "", err
		}
		tempFile, err := os.CreateTemp(filepath.Dir(backupPath), ".symlink-content-*")
		if err != nil {
			return "", err
		}
		tempPath := tempFile.Name()
		_, writeErr := tempFile.WriteString(target)
		closeErr := tempFile.Close()
		if writeErr != nil {
			_ = os.Remove(tempPath)
			return "", writeErr
		}
		if closeErr != nil {
			_ = os.Remove(tempPath)
			return "", closeErr
		}
		defer os.Remove(tempPath)
		hashPath = tempPath
	}
	return querySHA(backend, "hash-object", "--no-filters", "--", hashPath)
}

func listTrackedOpenChangeFiles(backend subshelldomain.RunnerQuerier) ([]string, error) {
	unstagedOutput, err := backend.QueryZ("git", "-c", "core.safecrlf=false", "diff", "--name-only", "-z")
	if err != nil {
		return []string{}, fmt.Errorf(messages.OpenChangesSnapshotCreateProblem, err)
	}
	stagedOutput, err := backend.QueryZ("git", "-c", "core.safecrlf=false", "diff", "--cached", "--name-only", "-z")
	if err != nil {
		return []string{}, fmt.Errorf(messages.OpenChangesSnapshotCreateProblem, err)
	}
	return mergeOpenChangesFileLists(zeroDelineatedFiles(unstagedOutput), zeroDelineatedFiles(stagedOutput)), nil
}

func listUntrackedFiles(backend subshelldomain.RunnerQuerier) ([]string, error) {
	output, err := backend.QueryZ("git", "ls-files", "--others", "--exclude-standard", "-z")
	if err != nil {
		return []string{}, fmt.Errorf(messages.OpenChangesSnapshotCreateProblem, err)
	}
	return zeroDelineatedFiles(output), nil
}

func mergeOpenChangesFileLists(first, second []string) []string {
	seen := make(map[string]bool, len(first)+len(second))
	result := make([]string, 0, len(first)+len(second))
	for _, file := range append(append([]string{}, first...), second...) {
		if len(file) == 0 || seen[file] {
			continue
		}
		seen[file] = true
		result = append(result, file)
	}
	sort.Strings(result)
	return result
}

func newOpenChangesFileMetadata(path string, info os.FileInfo) (openChangesFileMetadata, error) {
	result := openChangesFileMetadata{
		mode: openChangesFileMode(openChangesPermissionBits(info.Mode())),
		ownership: openChangesOwnership{
			groupID:    0,
			hasGroupID: false,
			hasUserID:  false,
			userID:     0,
		},
	}
	if result.mode&(os.ModeSetuid|os.ModeSetgid) == 0 {
		return result, nil
	}
	ownership, hasOwnership := openChangesFileOwnership(info)
	if !hasOwnership {
		return result, fmt.Errorf(messages.OpenChangesSnapshotOwnerProblem, path, errors.New("the operating system does not expose file ownership"))
	}
	if result.mode&os.ModeSetuid != 0 {
		result.ownership.userID = ownership.userID
		result.ownership.hasUserID = true
	}
	if result.mode&os.ModeSetgid != 0 {
		result.ownership.groupID = ownership.groupID
		result.ownership.hasGroupID = true
	}
	return result, nil
}

func openChangesBackupDir(gitDir stringss.Trimmed, snapshot gitdomain.SHA) string {
	return filepath.Join(gitDir.String(), filepath.FromSlash(openChangesBackupRoot), snapshot.String())
}

func openChangesBackupOwnershipMatches(info os.FileInfo, metadata openChangesFileMetadata) bool {
	if metadata.mode&(os.ModeSetuid|os.ModeSetgid) == 0 {
		return true
	}
	ownership, hasOwnership := openChangesFileOwnership(info)
	if !hasOwnership {
		return false
	}
	if metadata.mode&os.ModeSetuid != 0 && (!ownership.hasUserID || ownership.userID != metadata.ownership.userID) {
		return false
	}
	if metadata.mode&os.ModeSetgid != 0 && (!ownership.hasGroupID || ownership.groupID != metadata.ownership.groupID) {
		return false
	}
	return true
}

func openChangesBackupTypeMatches(info os.FileInfo, permission string) bool {
	switch permission {
	case "100644":
		return info.Mode().IsRegular() && info.Mode().Perm()&0o111 == 0
	case "100755":
		return info.Mode().IsRegular() && info.Mode().Perm()&0o111 != 0
	case "120000":
		return info.Mode()&os.ModeSymlink != 0
	default:
		return false
	}
}

func openChangesFileMode(permissionBits uint32) os.FileMode {
	result := os.FileMode(permissionBits & 0o777)
	if permissionBits&0o4000 != 0 {
		result |= os.ModeSetuid
	}
	if permissionBits&0o2000 != 0 {
		result |= os.ModeSetgid
	}
	if permissionBits&0o1000 != 0 {
		result |= os.ModeSticky
	}
	return result
}

func openChangesMetadataFile(value string) (string, bool) {
	fileBytes, err := hex.DecodeString(value)
	if err != nil || len(fileBytes) == 0 {
		return "", false
	}
	return string(fileBytes), true
}

func openChangesMetadataID(value string) (uint32, bool, bool) {
	if value == "-1" {
		return 0, false, true
	}
	id, err := strconv.ParseUint(value, 10, 32)
	return uint32(id), err == nil, err == nil
}

func openChangesMetadataIDText(id uint32, present bool) string {
	if !present {
		return "-1"
	}
	return strconv.FormatUint(uint64(id), 10)
}

func openChangesMetadataMatchEntries(entries map[string]openChangesTreeEntry, metadata map[string]openChangesFileMetadata) bool {
	regularFiles := 0
	for file, entry := range entries { // okay to iterate the map in random order because every entry is checked and order does not affect the result
		recordedMetadata, hasRecordedMetadata := metadata[file]
		if !openChangesRecordedModeMatchesEntry(recordedMetadata.mode, hasRecordedMetadata, entry.permission) {
			return false
		}
		if openChangesRegularPermission(entry.permission) {
			regularFiles++
		}
	}
	return regularFiles == len(metadata)
}

func openChangesMetadataMode(value string) (uint32, bool) {
	mode, err := strconv.ParseUint(value, 8, 32)
	return uint32(mode), err == nil && mode <= 0o7777
}

func openChangesOwnershipMatchesMode(mode os.FileMode, ownership openChangesOwnership) bool {
	if mode&os.ModeSetuid != 0 {
		if !ownership.hasUserID {
			return false
		}
	} else if ownership.hasUserID {
		return false
	}
	if mode&os.ModeSetgid != 0 {
		if !ownership.hasGroupID {
			return false
		}
	} else if ownership.hasGroupID {
		return false
	}
	return true
}

func openChangesPermissionBits(mode os.FileMode) uint32 {
	result := uint32(mode.Perm())
	if mode&os.ModeSetuid != 0 {
		result |= 0o4000
	}
	if mode&os.ModeSetgid != 0 {
		result |= 0o2000
	}
	if mode&os.ModeSticky != 0 {
		result |= 0o1000
	}
	return result
}

func openChangesRawFilesCommitMessage(modeLines []string) string {
	var result strings.Builder
	result.WriteString(openChangesRawFilesCommitTitle)
	result.WriteString("\n\n")
	result.WriteString(openChangesRawFileModesHeader)
	for _, line := range modeLines {
		result.WriteByte('\n')
		result.WriteString(line)
	}
	return result.String()
}

func openChangesRawFilesMetadata(querier subshelldomain.Querier, commit gitdomain.SHA) (map[string]openChangesFileMetadata, bool, error) {
	content, err := querier.Query("git", "cat-file", "commit", commit.String())
	if err != nil {
		return nil, false, err
	}
	_, message, found := strings.Cut(content, "\n\n")
	if !found {
		return nil, false, nil
	}
	subject, metadata, found := strings.Cut(message, "\n\n")
	if !found || subject != openChangesRawFilesCommitTitle {
		return nil, false, nil
	}
	lines := strings.Split(strings.TrimRight(metadata, "\n"), "\n")
	if len(lines) == 0 || lines[0] != openChangesRawFileModesHeader {
		return nil, false, nil
	}
	result := make(map[string]openChangesFileMetadata)
	for _, line := range lines[1:] {
		parts := strings.Split(line, " ")
		if len(parts) != 5 || parts[0] != "mode" || len(parts[1]) != 4 {
			return nil, false, nil
		}
		modeValue, validMode := openChangesMetadataMode(parts[1])
		userID, hasUserID, validUserID := openChangesMetadataID(parts[2])
		groupID, hasGroupID, validGroupID := openChangesMetadataID(parts[3])
		file, validFile := openChangesMetadataFile(parts[4])
		if !validMode || !validUserID || !validGroupID || !validFile {
			return nil, false, nil
		}
		if _, exists := result[file]; exists {
			return nil, false, nil
		}
		metadata := openChangesFileMetadata{
			mode: openChangesFileMode(modeValue),
			ownership: openChangesOwnership{
				groupID:    groupID,
				hasGroupID: hasGroupID,
				hasUserID:  hasUserID,
				userID:     userID,
			},
		}
		if !openChangesOwnershipMatchesMode(metadata.mode, metadata.ownership) {
			return nil, false, nil
		}
		result[file] = metadata
	}
	return result, true, nil
}

func openChangesRecordedModeMatchesEntry(recordedMode os.FileMode, hasRecordedMode bool, permission string) bool {
	switch permission {
	case "100644":
		return hasRecordedMode && recordedMode.Perm()&0o111 == 0
	case "100755":
		return hasRecordedMode && recordedMode.Perm()&0o111 != 0
	case "120000":
		return !hasRecordedMode
	default:
		return false
	}
}

func openChangesRegularPermission(permission string) bool {
	return permission == "100644" || permission == "100755"
}

func openChangesSnapshotRawCommit(querier subshelldomain.Querier, snapshot gitdomain.SHA) (Option[gitdomain.SHA], error) {
	content, err := querier.Query("git", "cat-file", "commit", snapshot.String())
	if err != nil {
		return None[gitdomain.SHA](), err
	}
	parents := make([]string, 0, 3)
	for line := range strings.Lines(content) {
		if parent, hasParent := strings.CutPrefix(line, "parent "); hasParent {
			parents = append(parents, strings.TrimSpace(parent))
		}
	}
	if len(parents) < 3 {
		return None[gitdomain.SHA](), nil
	}
	rawIndex := 2
	if len(parents) >= 4 {
		rawIndex = 3
	}
	rawFilesCommit, err := gitdomain.NewSHA(stringss.Trimmed(parents[rawIndex]))
	if err != nil {
		return None[gitdomain.SHA](), err
	}
	return Some(rawFilesCommit), nil
}

func openChangesTreeEntries(querier subshelldomain.Querier, commit gitdomain.SHA) (map[string]openChangesTreeEntry, error) {
	output, err := querier.QueryZ("git", "ls-tree", "-r", "-z", commit.String())
	if err != nil {
		return map[string]openChangesTreeEntry{}, err
	}
	result := make(map[string]openChangesTreeEntry)
	for _, entry := range output.Lines() {
		if len(entry) == 0 {
			continue
		}
		metadata, file, found := strings.Cut(entry, "\t")
		if !found {
			return map[string]openChangesTreeEntry{}, errors.New(messages.OpenChangesSnapshotInvalid)
		}
		metadataParts := strings.Fields(metadata)
		if len(metadataParts) != 3 {
			return map[string]openChangesTreeEntry{}, errors.New(messages.OpenChangesSnapshotInvalid)
		}
		sha, err := gitdomain.NewSHA(stringss.Trimmed(metadataParts[2]))
		if err != nil {
			return map[string]openChangesTreeEntry{}, err
		}
		result[file] = openChangesTreeEntry{
			objectType: metadataParts[1],
			permission: metadataParts[0],
			sha:        sha,
		}
	}
	return result, nil
}

func preserveSetIDOwnership(path string, metadata openChangesFileMetadata) error {
	return preserveSetIDOwnershipWith(path, metadata, openChangesLchown, openChangesFileOwnership)
}

func preserveSetIDOwnershipWith(path string, metadata openChangesFileMetadata, chown func(string, uint32, uint32, bool, bool) error, fileOwnership func(os.FileInfo) (openChangesOwnership, bool)) error {
	if metadata.mode&(os.ModeSetuid|os.ModeSetgid) == 0 {
		return nil
	}
	if !openChangesOwnershipMatchesMode(metadata.mode, metadata.ownership) {
		return fmt.Errorf(messages.OpenChangesSnapshotOwnerProblem, path, errors.New("the recorded ownership is incomplete"))
	}
	info, err := os.Lstat(path)
	if err != nil {
		return err
	}
	ownership, hasOwnership := fileOwnership(info)
	if !hasOwnership {
		return fmt.Errorf(messages.OpenChangesSnapshotOwnerProblem, path, errors.New("the operating system does not expose file ownership"))
	}
	changeUserID := false
	changeGroupID := false
	if metadata.mode&os.ModeSetuid != 0 && (!ownership.hasUserID || ownership.userID != metadata.ownership.userID) {
		changeUserID = true
	}
	if metadata.mode&os.ModeSetgid != 0 && (!ownership.hasGroupID || ownership.groupID != metadata.ownership.groupID) {
		changeGroupID = true
	}
	if changeUserID || changeGroupID {
		if err := chown(path, metadata.ownership.userID, metadata.ownership.groupID, changeUserID, changeGroupID); err != nil {
			return fmt.Errorf(messages.OpenChangesSnapshotOwnerProblem, path, err)
		}
	}
	info, err = os.Lstat(path)
	if err != nil {
		return err
	}
	ownership, hasOwnership = fileOwnership(info)
	if !hasOwnership {
		return fmt.Errorf(messages.OpenChangesSnapshotOwnerProblem, path, errors.New("the operating system does not expose file ownership"))
	}
	if metadata.mode&os.ModeSetuid != 0 && (!ownership.hasUserID || ownership.userID != metadata.ownership.userID) {
		return fmt.Errorf(messages.OpenChangesSnapshotOwnerProblem, path, errors.New("the operating system reports a different owner after chown"))
	}
	if metadata.mode&os.ModeSetgid != 0 && (!ownership.hasGroupID || ownership.groupID != metadata.ownership.groupID) {
		return fmt.Errorf(messages.OpenChangesSnapshotOwnerProblem, path, errors.New("the operating system reports a different group after chown"))
	}
	return nil
}

func querySHA(querier subshelldomain.Querier, args ...string) (gitdomain.SHA, error) {
	output, err := querier.QueryTrim("git", args...)
	if err != nil {
		return "", err
	}
	fields := strings.Fields(output.String())
	if len(fields) == 0 {
		return "", errors.New(messages.OpenChangesSnapshotInvalid)
	}
	return gitdomain.NewSHA(stringss.Trimmed(fields[len(fields)-1]))
}

func restoreFileExact(source, destination string, metadata openChangesFileMetadata) error {
	return restoreFileExactWith(source, destination, metadata, preserveSetIDOwnership)
}

func restoreFileExactWith(source, destination string, metadata openChangesFileMetadata, preserveOwnership func(string, openChangesFileMetadata) error) error {
	info, err := os.Lstat(source)
	if err != nil {
		return err
	}
	destinationDir := filepath.Dir(destination)
	if err := os.MkdirAll(destinationDir, 0o755); err != nil {
		return err
	}
	if info.Mode()&os.ModeSymlink != 0 {
		target, err := os.Readlink(source)
		if err != nil {
			return err
		}
		tempFile, err := os.CreateTemp(destinationDir, ".git-town-restore-*")
		if err != nil {
			return err
		}
		tempPath := tempFile.Name()
		if err := tempFile.Close(); err != nil {
			_ = os.Remove(tempPath)
			return err
		}
		if err := os.Remove(tempPath); err != nil {
			return err
		}
		if err := os.Symlink(target, tempPath); err != nil {
			return err
		}
		if err := os.Rename(tempPath, destination); err != nil {
			_ = os.Remove(tempPath)
			return err
		}
		return nil
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf(messages.OpenChangesSnapshotUnsupportedFile, source)
	}
	tempFile, err := os.CreateTemp(destinationDir, ".git-town-restore-*")
	if err != nil {
		return err
	}
	tempPath := tempFile.Name()
	if err := copyFileContents(source, tempFile); err != nil {
		_ = os.Remove(tempPath)
		return err
	}
	if err := preserveOwnership(tempPath, metadata); err != nil {
		_ = os.Remove(tempPath)
		return err
	}
	if err := chmodFileExact(tempPath, metadata.mode); err != nil {
		_ = os.Remove(tempPath)
		return err
	}
	if err := os.Rename(tempPath, destination); err != nil {
		_ = os.Remove(tempPath)
		return err
	}
	return nil
}

func sortedOpenChangesTreeFiles(entries map[string]openChangesTreeEntry) []string {
	result := make([]string, 0, len(entries))
	for file := range entries { // okay to iterate the map in random order because we sort the result below
		result = append(result, file)
	}
	sort.Strings(result)
	return result
}

func worktreeFilePath(rootDir gitdomain.RepoRootDir, file string) (string, error) {
	cleaned := path.Clean(file)
	if len(file) == 0 || path.IsAbs(file) || cleaned == ".." || strings.HasPrefix(cleaned, "../") {
		return "", fmt.Errorf(messages.OpenChangesSnapshotUnsupportedFile, file)
	}
	return filepath.Join(rootDir.String(), filepath.FromSlash(file)), nil
}

func zeroDelineatedFiles(output stringss.ZeroDelineated) []string {
	result := make([]string, 0)
	for _, file := range output.Lines() {
		if len(file) > 0 {
			result = append(result, file)
		}
	}
	return result
}
