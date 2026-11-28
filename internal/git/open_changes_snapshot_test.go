package git_test

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"testing"

	"github.com/git-town/git-town/v24/internal/git"
	"github.com/git-town/git-town/v24/internal/git/gitdomain"
	"github.com/git-town/git-town/v24/internal/gohacks/stringss"
	"github.com/git-town/git-town/v24/internal/subshell/subshelldomain"
	"github.com/git-town/git-town/v24/internal/test/commands"
	"github.com/git-town/git-town/v24/internal/test/testruntime"
	"github.com/git-town/git-town/v24/pkg/asserts"
	"github.com/shoenig/test/must"
)

func TestOpenChangesSnapshot(t *testing.T) {
	t.Parallel()

	t.Run("a failed raw restoration keeps the existing destination file", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		runtime.CreateFile("untracked", "original content")
		snapshot := createOpenChangesSnapshot(t, runtime)
		runtime.MustRun("git", "clean", "-fdx")
		runtime.CreateFile("untracked", "replacement content")
		backupPath := openChangesBackupFile(t, runtime, snapshot)
		asserts.NoError(os.Remove(backupPath))
		err := runtime.Git.RestoreOpenChangesSnapshot(runtime.TestRunner, runtime.TestRunner, snapshot)
		must.Error(t, err)
		content := asserts.NoError1(os.ReadFile(filepath.Join(runtime.WorkingDir, "untracked")))
		must.EqOp(t, "replacement content", string(content))
	})

	t.Run("no open changes", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		snapshot := asserts.NoError1(runtime.Git.CreateOpenChangesSnapshot(runtime.TestRunner))
		must.True(t, snapshot.IsNone())
	})

	t.Run("restores LF bytes under core.autocrlf=true", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		runtime.MustRun("git", "config", "core.autocrlf", "true")
		runtime.CreateFile(".gitattributes", "*.txt text\n")
		runtime.CreateFile("tracked.txt", "base content\n")
		runtime.MustRun("git", "add", ".gitattributes", "tracked.txt")
		runtime.MustRun("git", "commit", "-m", "base")
		worktreeBytes := []byte("worktree line 1\nworktree line 2\n")
		runtime.CreateFile("tracked.txt", string(worktreeBytes))
		snapshot := createOpenChangesSnapshot(t, runtime)
		runtime.MustRun("git", "reset", "--hard")
		asserts.NoError(runtime.Git.RestoreOpenChangesSnapshot(runtime.TestRunner, runtime.TestRunner, snapshot))
		restoredBytes := asserts.NoError1(os.ReadFile(filepath.Join(runtime.WorkingDir, "tracked.txt")))
		must.EqOp(t, string(worktreeBytes), string(restoredBytes))
	})

	t.Run("restores a snapshot containing only untracked files", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		runtime.CreateFile("untracked", "untracked content")
		snapshotOpt := asserts.NoError1(runtime.Git.CreateOpenChangesSnapshot(runtime.TestRunner))
		snapshot, hasSnapshot := snapshotOpt.Get()
		must.True(t, hasSnapshot)
		runtime.MustRun("git", "clean", "-fdx")
		asserts.NoError(runtime.Git.RestoreOpenChangesSnapshot(runtime.TestRunner, runtime.TestRunner, snapshot))
		status := runtime.MustQueryTrim("git", "status", "--porcelain=v1").String()
		must.EqOp(t, "?? untracked", status)
		content := asserts.NoError1(os.ReadFile(filepath.Join(runtime.WorkingDir, "untracked")))
		must.EqOp(t, "untracked content", string(content))
	})

	t.Run("restores recorded permission bits instead of current backup file mode", func(t *testing.T) {
		t.Parallel()
		if runtime.GOOS == "windows" {
			t.Skip("Windows does not preserve POSIX permission bits")
		}
		runtime := testruntime.Create(t)
		runtime.CreateFile("untracked", "untracked content")
		filePath := filepath.Join(runtime.WorkingDir, "untracked")
		asserts.NoError(os.Chmod(filePath, 0o644))
		snapshot := createOpenChangesSnapshot(t, runtime)
		backupPath := openChangesBackupFile(t, runtime, snapshot)
		asserts.NoError(os.Chmod(backupPath, 0o600))
		runtime.MustRun("git", "clean", "-fdx")
		asserts.NoError(runtime.Git.RestoreOpenChangesSnapshot(runtime.TestRunner, runtime.TestRunner, snapshot))
		info := asserts.NoError1(os.Lstat(filePath))
		must.EqOp(t, os.FileMode(0o644), info.Mode().Perm())
	})

	t.Run("restores recorded setgid ownership and permission bits", func(t *testing.T) {
		t.Parallel()
		if runtime.GOOS == "windows" {
			t.Skip("Windows does not preserve POSIX permission bits")
		}
		runtime := testruntime.Create(t)
		recordedGroup := openChangesAlternateGroup(t, runtime.WorkingDir)
		runtime.CreateFile("untracked", "untracked content")
		filePath := filepath.Join(runtime.WorkingDir, "untracked")
		asserts.NoError(os.Chown(filePath, -1, recordedGroup))
		asserts.NoError(os.Chmod(filePath, 0o644|os.ModeSetgid))
		originalInfo := asserts.NoError1(os.Lstat(filePath))
		must.True(t, originalInfo.Mode()&os.ModeSetgid != 0)
		snapshot := createOpenChangesSnapshot(t, runtime)
		backupPath := openChangesBackupFile(t, runtime, snapshot)
		_, backupGroup := openChangesTestFileOwnership(t, backupPath)
		must.EqOp(t, int64(recordedGroup), backupGroup)
		runtime.MustRun("git", "clean", "-fdx")
		asserts.NoError(runtime.Git.RestoreOpenChangesSnapshot(runtime.TestRunner, runtime.TestRunner, snapshot))
		info := asserts.NoError1(os.Lstat(filePath))
		must.EqOp(t, os.FileMode(0o644), info.Mode().Perm())
		must.True(t, info.Mode()&os.ModeSetgid != 0)
		_, restoredGroup := openChangesTestFileOwnership(t, filePath)
		must.EqOp(t, int64(recordedGroup), restoredGroup)
	})

	t.Run("restores recorded special permission bits", func(t *testing.T) {
		t.Parallel()
		if runtime.GOOS == "windows" {
			t.Skip("Windows does not preserve POSIX permission bits")
		}
		runtime := testruntime.Create(t)
		runtime.CreateFile("untracked", "untracked content")
		filePath := filepath.Join(runtime.WorkingDir, "untracked")
		asserts.NoError(os.Chmod(filePath, 0o644|os.ModeSticky))
		snapshot := createOpenChangesSnapshot(t, runtime)
		runtime.MustRun("git", "clean", "-fdx")
		asserts.NoError(runtime.Git.RestoreOpenChangesSnapshot(runtime.TestRunner, runtime.TestRunner, snapshot))
		info := asserts.NoError1(os.Lstat(filePath))
		must.EqOp(t, os.FileMode(0o644), info.Mode().Perm())
		must.True(t, info.Mode()&os.ModeSticky != 0)
	})

	t.Run("restores staged, unstaged, deleted, and untracked files exactly", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		runtime.CreateFile(".gitattributes", "*.txt text eol=lf")
		runtime.CreateFile("deleted", "deleted content")
		runtime.CreateFile("tracked.txt", "base content\n")
		runtime.MustRun("git", "add", ".gitattributes", "deleted", "tracked.txt")
		runtime.MustRun("git", "commit", "-m", "base")
		stagedBytes := []byte("staged line 1\r\nstaged line 2\r\n")
		runtime.CreateFile("tracked.txt", string(stagedBytes))
		runtime.MustRun("git", "add", "tracked.txt")
		worktreeBytes := []byte("worktree line 1\r\nworktree line 2\r\n")
		runtime.CreateFile("tracked.txt", string(worktreeBytes))
		runtime.MustRun("git", "rm", "deleted")
		untrackedName := ":(exclude)file.txt"
		untrackedBytes := []byte("line 1\r\nline 2\r\n")
		runtime.CreateFile(untrackedName, string(untrackedBytes))
		beforeStatus := runtime.MustQueryTrim("git", "status", "--porcelain=v1").String()
		snapshotOpt := asserts.NoError1(runtime.Git.CreateOpenChangesSnapshot(runtime.TestRunner))
		snapshot, hasSnapshot := snapshotOpt.Get()
		must.True(t, hasSnapshot)
		snapshotRef := git.OpenChangesSnapshotRef(snapshot)
		refSHA := runtime.MustQueryTrim("git", "rev-parse", "--verify", "-q", snapshotRef).String()
		must.EqOp(t, snapshot.String(), refSHA)
		asserts.NoError(runtime.Git.ValidateOpenChangesSnapshot(runtime.TestRunner, snapshot))
		runtime.MustRun("git", "gc", "--prune=now")
		runtime.MustRun("git", "reset", "--hard")
		runtime.MustRun("git", "clean", "-fdx")
		asserts.NoError(runtime.Git.RestoreOpenChangesSnapshot(runtime.TestRunner, runtime.TestRunner, snapshot))
		afterStatus := runtime.MustQueryTrim("git", "status", "--porcelain=v1").String()
		must.EqOp(t, beforeStatus, afterStatus)
		trackedContent := asserts.NoError1(os.ReadFile(filepath.Join(runtime.WorkingDir, "tracked.txt")))
		must.EqOp(t, string(worktreeBytes), string(trackedContent))
		normalizedIndexBytes := []byte("staged line 1\nstaged line 2\n")
		normalizedIndexPath := filepath.Join(t.TempDir(), "normalized-index")
		asserts.NoError(os.WriteFile(normalizedIndexPath, normalizedIndexBytes, 0o600))
		expectedIndexSHA := runtime.MustQueryTrim("git", "hash-object", "--no-filters", "--", normalizedIndexPath).String()
		indexSHA := runtime.MustQueryTrim("git", "rev-parse", ":0:tracked.txt").String()
		must.EqOp(t, expectedIndexSHA, indexSHA)
		restoredUntrackedBytes := asserts.NoError1(os.ReadFile(filepath.Join(runtime.WorkingDir, untrackedName)))
		must.EqOp(t, string(untrackedBytes), string(restoredUntrackedBytes))
		refResult := runtime.MustQueryStringCode("git rev-parse --verify -q " + snapshotRef)
		must.EqOp(t, 1, refResult.ExitCode)
	})

	t.Run("stores a large mode manifest through a commit-tree file", func(t *testing.T) {
		t.Parallel()
		if runtime.GOOS == "windows" {
			t.Skip("uses paths longer than Windows supports in test repositories")
		}
		runtime := testruntime.Create(t)
		longDir := strings.Repeat("x", 240)
		longSubdir := strings.Repeat("y", 240)
		for index := range 130 {
			runtime.CreateFile(fmt.Sprintf("large-file-%03d/%s/%s/content", index, longDir, longSubdir), "x")
		}
		runner := &gitRecordingRunner{backend: runtime.TestRunner}
		snapshotOpt := asserts.NoError1(runtime.Git.CreateOpenChangesSnapshot(runner))
		snapshot, hasSnapshot := snapshotOpt.Get()
		must.True(t, hasSnapshot)
		fileUses := 0
		for _, command := range runner.commands {
			if len(command) < 2 || command[0] != "git" || command[1] != "commit-tree" {
				continue
			}
			for _, arg := range command[2:] {
				must.True(t, len(arg) < 1024)
			}
			if slices.Contains(command, "-F") {
				fileUses++
			}
		}
		must.EqOp(t, 1, fileUses)
		must.True(t, len(runner.commitTreeFileContent) > 128*1024)
		must.EqOp(t, 130, strings.Count(runner.commitTreeFileContent, "\nmode "))
		_, err := os.Stat(runner.commitTreeMessagePath)
		must.True(t, os.IsNotExist(err))
		rawCommit := runtime.MustQueryTrim("git", "rev-parse", snapshot.String()+"^4").String()
		rawCommitObject := runtime.MustQuery("git", "cat-file", "commit", rawCommit)
		_, rawCommitMessage, found := strings.Cut(rawCommitObject, "\n\n")
		must.True(t, found)
		must.EqOp(t, runner.commitTreeFileContent, rawCommitMessage)
	})

	t.Run("validation rejects a changed raw backup file", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		runtime.CreateFile("untracked", "untracked content")
		snapshot := createOpenChangesSnapshot(t, runtime)
		backupPath := openChangesBackupFile(t, runtime, snapshot)
		asserts.NoError(os.WriteFile(backupPath, []byte("changed content"), 0o600))
		err := runtime.Git.ValidateOpenChangesSnapshot(runtime.TestRunner, snapshot)
		must.Error(t, err)
		must.EqOp(t, "the saved open changes snapshot is incomplete", err.Error())
	})

	t.Run("validation rejects a missing raw backup file", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		runtime.CreateFile("untracked", "untracked content")
		snapshot := createOpenChangesSnapshot(t, runtime)
		backupPath := openChangesBackupFile(t, runtime, snapshot)
		asserts.NoError(os.Remove(backupPath))
		err := runtime.Git.ValidateOpenChangesSnapshot(runtime.TestRunner, snapshot)
		must.Error(t, err)
		must.EqOp(t, "the saved open changes snapshot is incomplete", err.Error())
	})

	t.Run("validation rejects a raw backup file with a changed sticky bit", func(t *testing.T) {
		t.Parallel()
		if runtime.GOOS == "windows" {
			t.Skip("Windows does not preserve POSIX permission bits")
		}
		runtime := testruntime.Create(t)
		runtime.CreateFile("untracked", "untracked content")
		filePath := filepath.Join(runtime.WorkingDir, "untracked")
		asserts.NoError(os.Chmod(filePath, 0o644))
		snapshot := createOpenChangesSnapshot(t, runtime)
		backupPath := openChangesBackupFile(t, runtime, snapshot)
		asserts.NoError(os.Chmod(backupPath, 0o644|os.ModeSticky))
		err := runtime.Git.ValidateOpenChangesSnapshot(runtime.TestRunner, snapshot)
		must.Error(t, err)
		must.EqOp(t, "the saved open changes snapshot is incomplete", err.Error())
	})

	t.Run("validation rejects a raw backup file with changed permission bits", func(t *testing.T) {
		t.Parallel()
		if runtime.GOOS == "windows" {
			t.Skip("Windows does not preserve POSIX permission bits")
		}
		runtime := testruntime.Create(t)
		runtime.CreateFile("untracked", "untracked content")
		filePath := filepath.Join(runtime.WorkingDir, "untracked")
		asserts.NoError(os.Chmod(filePath, 0o644))
		snapshot := createOpenChangesSnapshot(t, runtime)
		backupPath := openChangesBackupFile(t, runtime, snapshot)
		asserts.NoError(os.Chmod(backupPath, 0o600))
		err := runtime.Git.ValidateOpenChangesSnapshot(runtime.TestRunner, snapshot)
		must.Error(t, err)
		must.EqOp(t, "the saved open changes snapshot is incomplete", err.Error())
	})

	t.Run("validation rejects a raw backup file with changed setgid ownership", func(t *testing.T) {
		t.Parallel()
		if runtime.GOOS == "windows" {
			t.Skip("Windows does not preserve POSIX permission bits")
		}
		groups := asserts.NoError1(os.Getgroups())
		otherGroup := -1
		for _, group := range groups {
			if group != os.Getegid() {
				otherGroup = group
				break
			}
		}
		if otherGroup == -1 {
			t.Skip("test user needs a second group")
		}
		runtime := testruntime.Create(t)
		runtime.CreateFile("untracked", "untracked content")
		filePath := filepath.Join(runtime.WorkingDir, "untracked")
		asserts.NoError(os.Chown(filePath, -1, os.Getegid()))
		asserts.NoError(os.Chmod(filePath, 0o644|os.ModeSetgid))
		snapshot := createOpenChangesSnapshot(t, runtime)
		backupPath := openChangesBackupFile(t, runtime, snapshot)
		_, originalBackupGroup := openChangesTestFileOwnership(t, backupPath)
		must.EqOp(t, int64(os.Getegid()), originalBackupGroup)
		asserts.NoError(os.Chown(backupPath, -1, otherGroup))
		asserts.NoError(os.Chmod(backupPath, 0o644|os.ModeSetgid))
		_, tamperedBackupGroup := openChangesTestFileOwnership(t, backupPath)
		must.EqOp(t, int64(otherGroup), tamperedBackupGroup)
		backupInfo := asserts.NoError1(os.Lstat(backupPath))
		must.EqOp(t, os.FileMode(0o644), backupInfo.Mode().Perm())
		must.True(t, backupInfo.Mode()&os.ModeSetgid != 0)
		err := runtime.Git.ValidateOpenChangesSnapshot(runtime.TestRunner, snapshot)
		must.Error(t, err)
		must.EqOp(t, "the saved open changes snapshot is incomplete", err.Error())
	})

	t.Run("validation rejects a raw backup file with the wrong type", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		runtime.CreateFile("untracked", "untracked content")
		snapshot := createOpenChangesSnapshot(t, runtime)
		backupPath := openChangesBackupFile(t, runtime, snapshot)
		asserts.NoError(os.Remove(backupPath))
		asserts.NoError(os.Mkdir(backupPath, 0o755))
		err := runtime.Git.ValidateOpenChangesSnapshot(runtime.TestRunner, snapshot)
		must.Error(t, err)
		must.EqOp(t, "the saved open changes snapshot is incomplete", err.Error())
	})
}

func createOpenChangesSnapshot(t *testing.T, runtime commands.TestCommands) gitdomain.SHA {
	t.Helper()
	snapshotOpt := asserts.NoError1(runtime.Git.CreateOpenChangesSnapshot(runtime.TestRunner))
	snapshot, hasSnapshot := snapshotOpt.Get()
	must.True(t, hasSnapshot)
	return snapshot
}

func openChangesAlternateGroup(t *testing.T, directory string) int {
	t.Helper()
	probePath := filepath.Join(directory, ".group-probe")
	asserts.NoError(os.WriteFile(probePath, []byte("probe"), 0o600))
	_, probeGroup := openChangesTestFileOwnership(t, probePath)
	asserts.NoError(os.Remove(probePath))
	groups := asserts.NoError1(os.Getgroups())
	groups = append(groups, os.Getegid())
	for _, group := range groups {
		if group >= 0 && int64(group) != probeGroup {
			return group
		}
	}
	t.Skip("test user needs a group different from newly created files")
	return -1
}

func openChangesBackupFile(t *testing.T, runtime commands.TestCommands, snapshot gitdomain.SHA) string {
	t.Helper()
	gitDir := runtime.MustQueryTrim("git", "rev-parse", "--absolute-git-dir").String()
	return filepath.Join(gitDir, "git-town", "open-changes", snapshot.String(), "files", "untracked")
}

type gitRecordingRunner struct {
	backend               subshelldomain.RunnerQuerier
	commands              [][]string
	commitTreeFileContent string
	commitTreeMessagePath string
}

func (self *gitRecordingRunner) Query(executable string, args ...string) (string, error) {
	self.record(executable, args...)
	return self.backend.Query(executable, args...)
}

func (self *gitRecordingRunner) QueryTrim(executable string, args ...string) (stringss.Trimmed, error) {
	self.record(executable, args...)
	if executable == "git" && len(args) > 0 && args[0] == "commit-tree" {
		if index := slices.Index(args, "-F"); index >= 0 && index+1 < len(args) {
			self.commitTreeMessagePath = args[index+1]
			content, err := os.ReadFile(self.commitTreeMessagePath)
			if err == nil {
				self.commitTreeFileContent = string(content)
			}
		}
	}
	return self.backend.QueryTrim(executable, args...)
}

func (self *gitRecordingRunner) QueryWithEnv(env []string, executable string, args ...string) (string, error) {
	self.record(executable, args...)
	return self.backend.QueryWithEnv(env, executable, args...)
}

func (self *gitRecordingRunner) QueryZ(executable string, args ...string) (stringss.ZeroDelineated, error) {
	self.record(executable, args...)
	return self.backend.QueryZ(executable, args...)
}

func (self *gitRecordingRunner) Run(executable string, args ...string) error {
	self.record(executable, args...)
	return self.backend.Run(executable, args...)
}

func (self *gitRecordingRunner) RunWithEnv(env []string, executable string, args ...string) error {
	self.record(executable, args...)
	return self.backend.RunWithEnv(env, executable, args...)
}

func (self *gitRecordingRunner) record(executable string, args ...string) {
	command := append([]string{executable}, args...)
	self.commands = append(self.commands, command)
}
