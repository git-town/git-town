package git_test

import (
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/git-town/git-town/src/git"
	"github.com/git-town/git-town/test"
	"github.com/stretchr/testify/assert"
)

func TestRunner_AddRemote(t *testing.T) {
	runner := CreateTestGitTownRepo(t).Runner
	err := runner.AddRemote("foo", "bar")
	assert.Nil(t, err)
	remotes, err := runner.Remotes()
	assert.Nil(t, err)
	assert.Equal(t, []string{"foo"}, remotes)
}

func TestRunner_CheckoutBranch(t *testing.T) {
	runner := test.CreateRepo(t).Runner
	err := runner.CreateBranch("branch1", "master")
	assert.Nil(t, err)
	err = runner.CheckoutBranch("branch1")
	assert.Nil(t, err)
	currentBranch, err := runner.CurrentBranch()
	assert.Nil(t, err)
	assert.Equal(t, "branch1", currentBranch)
	err = runner.CheckoutBranch("master")
	assert.Nil(t, err)
	currentBranch, err = runner.CurrentBranch()
	assert.Nil(t, err)
	assert.Equal(t, "master", currentBranch)
}

func TestRunner_Commits(t *testing.T) {
	runner := test.CreateRepo(t).Runner
	err := runner.CreateCommit(git.Commit{
		Branch:      "master",
		FileName:    "file1",
		FileContent: "hello",
		Message:     "first commit",
	})
	assert.Nil(t, err)
	err = runner.CreateCommit(git.Commit{
		Branch:      "master",
		FileName:    "file2",
		FileContent: "hello again",
		Message:     "second commit",
	})
	assert.Nil(t, err)
	commits, err := runner.Commits([]string{"FILE NAME", "FILE CONTENT"})
	assert.Nil(t, err)
	assert.Len(t, commits, 2)
	assert.Equal(t, "master", commits[0].Branch)
	assert.Equal(t, "file1", commits[0].FileName)
	assert.Equal(t, "hello", commits[0].FileContent)
	assert.Equal(t, "first commit", commits[0].Message)
	assert.Equal(t, "master", commits[1].Branch)
	assert.Equal(t, "file2", commits[1].FileName)
	assert.Equal(t, "hello again", commits[1].FileContent)
	assert.Equal(t, "second commit", commits[1].Message)
}

func TestRunner_Configuration(t *testing.T) {
	runner := test.CreateRepo(t).Runner
	config := runner.Configuration
	assert.NotNil(t, config, "first path: new config")
	config = runner.Configuration
	assert.NotNil(t, config, "second path: cached config")
}

func TestRunner_ConnectTrackingBranch(t *testing.T) {
	// replicating the situation this is used in,
	// connecting branches of repos with the same commits in them
	origin := test.CreateRepo(t)
	repoDir := filepath.Join(test.CreateTempDir(t), "repo") // need a non-existing directory
	err := test.CopyDirectory(origin.WorkingDir(), repoDir)
	assert.Nil(t, err)
	runner := test.NewRepo(repoDir, repoDir, "").Runner
	err = runner.AddRemote("origin", origin.WorkingDir())
	assert.Nil(t, err)
	err = runner.Fetch()
	assert.Nil(t, err)
	err = runner.ConnectTrackingBranch("master")
	assert.Nil(t, err)
	err = runner.PushBranch("master")
	assert.Nil(t, err)
}

func TestRunner_CreateBranch(t *testing.T) {
	runner := test.CreateRepo(t).Runner
	err := runner.CreateBranch("branch1", "master")
	assert.Nil(t, err)
	currentBranch, err := runner.CurrentBranch()
	assert.Nil(t, err)
	assert.Equal(t, "master", currentBranch)
	branches, err := runner.LocalBranches()
	assert.Nil(t, err)
	assert.Equal(t, []string{"branch1", "master"}, branches)
}

func TestRunner_CreateChildFeatureBranch(t *testing.T) {
	runner := CreateTestGitTownRepo(t).Runner
	err := runner.CreateFeatureBranch("f1")
	assert.Nil(t, err)
	err = runner.CreateChildFeatureBranch("f1a", "f1")
	assert.Nil(t, err)
	res, err := runner.Run("git", "town", "config")
	assert.Nil(t, err)
	has := strings.Contains(res.OutputSanitized(), "Branch Ancestry:\n  main\n    f1\n      f1a")
	assert.True(t, has)
}

func TestRunner_CreateCommit(t *testing.T) {
	runner := test.CreateRepo(t).Runner
	err := runner.CreateCommit(git.Commit{
		Branch:      "master",
		FileName:    "hello.txt",
		FileContent: "hello world",
		Message:     "test commit",
	})
	assert.Nil(t, err)
	commits, err := runner.Commits([]string{"FILE NAME", "FILE CONTENT"})
	assert.Nil(t, err)
	assert.Len(t, commits, 1)
	assert.Equal(t, "hello.txt", commits[0].FileName)
	assert.Equal(t, "hello world", commits[0].FileContent)
	assert.Equal(t, "test commit", commits[0].Message)
	assert.Equal(t, "master", commits[0].Branch)
}

func TestRunner_CreateCommit_Author(t *testing.T) {
	runner := test.CreateRepo(t).Runner
	err := runner.CreateCommit(git.Commit{
		Branch:      "master",
		FileName:    "hello.txt",
		FileContent: "hello world",
		Message:     "test commit",
		Author:      "developer <developer@example.com>",
	})
	assert.Nil(t, err)
	commits, err := runner.Commits([]string{"FILE NAME", "FILE CONTENT"})
	assert.Nil(t, err)
	assert.Len(t, commits, 1)
	assert.Equal(t, "hello.txt", commits[0].FileName)
	assert.Equal(t, "hello world", commits[0].FileContent)
	assert.Equal(t, "test commit", commits[0].Message)
	assert.Equal(t, "master", commits[0].Branch)
	assert.Equal(t, "developer <developer@example.com>", commits[0].Author)
}

func TestRunner_CreateFeatureBranch(t *testing.T) {
	runner := CreateTestGitTownRepo(t).Runner
	err := runner.CreateFeatureBranch("f1")
	assert.Nil(t, err)
	runner.Configuration.Reload()
	assert.True(t, runner.Configuration.IsFeatureBranch("f1"))
	assert.Equal(t, []string{"main"}, runner.Configuration.GetAncestorBranches("f1"))
}

func TestRunner_CreateFeatureBranchNoParent(t *testing.T) {
	runner := CreateTestGitTownRepo(t).Runner
	err := runner.CreateFeatureBranchNoParent("f1")
	assert.Nil(t, err)
	runner.Configuration.Reload()
	assert.True(t, runner.Configuration.IsFeatureBranch("f1"))
	assert.Equal(t, []string(nil), runner.Configuration.GetAncestorBranches("f1"))
}

func TestRunner_CreateFile(t *testing.T) {
	runner := test.CreateRepo(t).Runner
	err := runner.CreateFile("filename", "content")
	assert.Nil(t, err, "cannot create file in repo")
	content, err := ioutil.ReadFile(filepath.Join(runner.WorkingDir(), "filename"))
	assert.Nil(t, err, "cannot read file")
	assert.Equal(t, "content", string(content))
}

func TestRunner_CreateFile_InSubFolder(t *testing.T) {
	runner := test.CreateRepo(t).Runner
	err := runner.CreateFile("folder/filename", "content")
	assert.Nil(t, err, "cannot create file in repo")
	content, err := ioutil.ReadFile(filepath.Join(runner.WorkingDir(), "folder/filename"))
	assert.Nil(t, err, "cannot read file")
	assert.Equal(t, "content", string(content))
}

func TestRunner_CreatePerennialBranches(t *testing.T) {
	runner := CreateTestGitTownRepo(t).Runner
	err := runner.CreatePerennialBranches("p1", "p2")
	assert.Nil(t, err)
	branches, err := runner.LocalBranches()
	assert.Nil(t, err)
	assert.Equal(t, []string{"main", "master", "p1", "p2"}, branches)
	runner.Configuration.Reload()
	assert.True(t, runner.Configuration.IsPerennialBranch("p1"))
	assert.True(t, runner.Configuration.IsPerennialBranch("p2"))
}

func TestRunner_CurrentBranch(t *testing.T) {
	runner := test.CreateRepo(t).Runner
	err := runner.CheckoutBranch("master")
	assert.Nil(t, err)
	err = runner.CreateBranch("b1", "master")
	assert.Nil(t, err)
	err = runner.CheckoutBranch("b1")
	assert.Nil(t, err)
	branch, err := runner.CurrentBranch()
	assert.Nil(t, err)
	assert.Equal(t, "b1", branch)
	err = runner.CheckoutBranch("master")
	assert.Nil(t, err)
	branch, err = runner.CurrentBranch()
	assert.Nil(t, err)
	assert.Equal(t, "master", branch)
}

func TestRunner_Fetch(t *testing.T) {
	runner := test.CreateRepo(t).Runner
	origin := test.CreateRepo(t)
	err := runner.AddRemote("origin", origin.WorkingDir())
	assert.Nil(t, err)
	err = runner.Fetch()
	assert.Nil(t, err)
}

func TestRunner_FileContentInCommit(t *testing.T) {
	runner := test.CreateRepo(t).Runner
	err := runner.CreateCommit(git.Commit{
		Branch:      "master",
		FileName:    "hello.txt",
		FileContent: "hello world",
		Message:     "commit",
	})
	assert.Nil(t, err)
	commits, err := runner.CommitsInBranch("master", []string{})
	assert.Nil(t, err)
	assert.Len(t, commits, 1)
	content, err := runner.FileContentInCommit(commits[0].SHA, "hello.txt")
	assert.Nil(t, err)
	assert.Equal(t, "hello world", content)
}

func TestRunner_FilesInCommit(t *testing.T) {
	runner := test.CreateRepo(t).Runner
	err := runner.CreateFile("f1.txt", "one")
	assert.Nil(t, err)
	err = runner.CreateFile("f2.txt", "two")
	assert.Nil(t, err)
	err = runner.StageFiles("f1.txt", "f2.txt")
	assert.Nil(t, err)
	err = runner.CommitStagedChanges("stuff")
	assert.Nil(t, err)
	commits, err := runner.Commits([]string{})
	assert.Nil(t, err)
	assert.Len(t, commits, 1)
	fileNames, err := runner.FilesInCommit(commits[0].SHA)
	assert.Nil(t, err)
	assert.Equal(t, []string{"f1.txt", "f2.txt"}, fileNames)
}

func TestRunner_HasBranchesOutOfSync_synced(t *testing.T) {
	env, err := test.NewStandardGitEnvironment(test.CreateTempDir(t))
	assert.Nil(t, err)
	runner := env.DevRepo.Runner
	err = runner.CreateBranch("branch1", "main")
	assert.Nil(t, err)
	err = runner.CheckoutBranch("branch1")
	assert.Nil(t, err)
	err = runner.CreateFile("file1", "content")
	assert.Nil(t, err)
	err = runner.StageFiles("file1")
	assert.Nil(t, err)
	err = runner.CommitStagedChanges("stuff")
	assert.Nil(t, err)
	err = runner.PushBranch("main")
	assert.Nil(t, err)
	have, err := runner.HasBranchesOutOfSync()
	assert.Nil(t, err)
	assert.False(t, have)
}

func TestRunner_HasBranchesOutOfSync_branchAhead(t *testing.T) {
	env, err := test.NewStandardGitEnvironment(test.CreateTempDir(t))
	assert.Nil(t, err)
	runner := env.DevRepo.Runner
	err = runner.CreateBranch("branch1", "main")
	assert.Nil(t, err)
	err = runner.PushBranch("branch1")
	assert.Nil(t, err)
	err = runner.CreateFile("file1", "content")
	assert.Nil(t, err)
	err = runner.StageFiles("file1")
	assert.Nil(t, err)
	err = runner.CommitStagedChanges("stuff")
	assert.Nil(t, err)
	have, err := runner.HasBranchesOutOfSync()
	assert.Nil(t, err)
	assert.True(t, have)
}

func TestRunner_HasBranchesOutOfSync_branchBehind(t *testing.T) {
	env, err := test.NewStandardGitEnvironment(test.CreateTempDir(t))
	assert.Nil(t, err)
	err = env.DevRepo.CreateBranch("branch1", "main")
	assert.Nil(t, err)
	err = env.DevRepo.PushBranch("branch1")
	assert.Nil(t, err)
	err = env.OriginRepo.CheckoutBranch("main")
	assert.Nil(t, err)
	err = env.OriginRepo.CreateFile("file1", "content")
	assert.Nil(t, err)
	err = env.OriginRepo.StageFiles("file1")
	assert.Nil(t, err)
	err = env.OriginRepo.CommitStagedChanges("stuff")
	assert.Nil(t, err)
	err = env.OriginRepo.CheckoutBranch("master")
	assert.Nil(t, err)
	err = env.DevRepo.Fetch()
	assert.Nil(t, err)
	have, err := env.DevRepo.Runner.HasBranchesOutOfSync()
	assert.Nil(t, err)
	assert.True(t, have)
}

func TestRunner_HasGitTownConfigNow(t *testing.T) {
	runner := test.CreateRepo(t).Runner
	res, err := runner.HasGitTownConfigNow()
	assert.Nil(t, err)
	assert.False(t, res)
	err = runner.CreateBranch("main", "master")
	assert.Nil(t, err)
	err = runner.CreateFeatureBranch("foo")
	assert.Nil(t, err)
	res, err = runner.HasGitTownConfigNow()
	assert.Nil(t, err)
	assert.True(t, res)
}

func TestRunner_HasFile(t *testing.T) {
	runner := test.CreateRepo(t).Runner
	err := runner.CreateFile("f1.txt", "one")
	assert.Nil(t, err)
	has, err := runner.HasFile("f1.txt", "one")
	assert.Nil(t, err)
	assert.True(t, has)
	_, err = runner.HasFile("f1.txt", "zonk")
	assert.Error(t, err)
	_, err = runner.HasFile("zonk.txt", "one")
	assert.Error(t, err)
}
func TestRunner_HasLocalBranch(t *testing.T) {
	origin := test.CreateRepo(t)
	repoDir := test.CreateTempDir(t)
	repo, err := origin.Clone(repoDir)
	assert.Nil(t, err)
	err = repo.CreateBranch("b1", "master")
	assert.Nil(t, err)
	err = repo.CreateBranch("b2", "master")
	assert.Nil(t, err)
	has, err := repo.HasLocalBranch("b1")
	assert.Nil(t, err)
	assert.True(t, has)
	has, err = repo.HasLocalBranch("b2")
	assert.Nil(t, err)
	assert.True(t, has)
	has, err = repo.HasLocalBranch("b3")
	assert.Nil(t, err)
	assert.False(t, has)
}

func TestRunner_HasOpenChanges(t *testing.T) {
	runner := test.CreateRepo(t).Runner
	has, err := runner.HasOpenChanges()
	assert.Nil(t, err)
	assert.False(t, has)
	err = runner.CreateFile("foo", "bar")
	assert.Nil(t, err)
	has, err = runner.HasOpenChanges()
	assert.Nil(t, err)
	assert.True(t, has)
}

func TestRunner_HasRebaseInProgress(t *testing.T) {
	runner := test.CreateRepo(t).Runner
	has, err := runner.HasRebaseInProgress()
	assert.Nil(t, err)
	assert.False(t, has)
}

func TestRunner_HasRemote(t *testing.T) {
	origin := test.CreateRepo(t)
	repoDir := test.CreateTempDir(t)
	repo, err := origin.Clone(repoDir)
	assert.Nil(t, err)
	has, err := repo.Runner.HasRemote("origin")
	assert.Nil(t, err)
	assert.True(t, has)
	has, err = repo.Runner.HasRemote("zonk")
	assert.Nil(t, err)
	assert.False(t, has)
}

func TestRunner_HasTrackingBranch(t *testing.T) {
	origin := test.CreateRepo(t)
	err := origin.CreateBranch("b1", "master")
	assert.Nil(t, err)
	repoDir := test.CreateTempDir(t)
	repo, err := origin.Clone(repoDir)
	assert.Nil(t, err)
	runner := repo.Runner
	err = runner.CheckoutBranch("b1")
	assert.Nil(t, err)
	err = runner.CreateBranch("b2", "master")
	assert.Nil(t, err)
	has, err := runner.HasTrackingBranch("b1")
	assert.Nil(t, err)
	assert.True(t, has)
	has, err = runner.HasTrackingBranch("b2")
	assert.Nil(t, err)
	assert.False(t, has)
	has, err = runner.HasTrackingBranch("b3")
	assert.Nil(t, err)
	assert.False(t, has)
}

func TestRunner_LastActiveDir(t *testing.T) {
	runner := test.CreateRepo(t).Runner
	dir, err := runner.LastActiveDir()
	assert.Nil(t, err)
	assert.Equal(t, runner.WorkingDir(), dir)
}

func TestRunner_LocalBranches(t *testing.T) {
	origin := test.CreateRepo(t)
	repoDir := test.CreateTempDir(t)
	repo, err := origin.Clone(repoDir)
	assert.Nil(t, err)
	err = repo.CreateBranch("b1", "master")
	assert.Nil(t, err)
	err = repo.CreateBranch("b2", "master")
	assert.Nil(t, err)
	err = origin.CreateBranch("b3", "master")
	assert.Nil(t, err)
	err = repo.Fetch()
	assert.Nil(t, err)
	branches, err := repo.Runner.LocalBranches()
	assert.Nil(t, err)
	assert.Equal(t, []string{"b1", "b2", "master"}, branches)
}

func TestRunner_LocalAndRemoteBranches(t *testing.T) {
	origin := test.CreateRepo(t)
	repoDir := test.CreateTempDir(t)
	repo, err := origin.Clone(repoDir)
	assert.Nil(t, err)
	err = repo.CreateBranch("b1", "master")
	assert.Nil(t, err)
	err = repo.CreateBranch("b2", "master")
	assert.Nil(t, err)
	err = origin.CreateBranch("b3", "master")
	assert.Nil(t, err)
	err = repo.Fetch()
	assert.Nil(t, err)
	branches, err := repo.Runner.LocalAndRemoteBranches()
	assert.Nil(t, err)
	assert.Equal(t, []string{"b1", "b2", "b3", "master"}, branches)
}

func TestRunner_PreviouslyCheckedOutBranch(t *testing.T) {
	runner := test.CreateRepo(t).Runner
	err := runner.CreateBranch("feature1", "master")
	assert.Nil(t, err)
	err = runner.CreateBranch("feature2", "master")
	assert.Nil(t, err)
	err = runner.CheckoutBranch("feature1")
	assert.Nil(t, err)
	err = runner.CheckoutBranch("feature2")
	assert.Nil(t, err)
	have, err := runner.PreviouslyCheckedOutBranch()
	assert.Nil(t, err)
	assert.Equal(t, "feature1", have)
}

func TestRunner_PushBranch(t *testing.T) {
	runner := test.CreateRepo(t).Runner
	origin := test.CreateRepo(t)
	err := runner.AddRemote("origin", origin.WorkingDir())
	assert.Nil(t, err)
	err = runner.CreateBranch("b1", "master")
	assert.Nil(t, err)
	err = runner.PushBranch("b1")
	assert.Nil(t, err)
	branches, err := origin.LocalBranches()
	assert.Nil(t, err)
	assert.Equal(t, []string{"b1", "master"}, branches)
}

func TestRunner_RemoteBranches(t *testing.T) {
	origin := test.CreateRepo(t)
	repoDir := test.CreateTempDir(t)
	repo, err := origin.Clone(repoDir)
	assert.Nil(t, err)
	err = repo.CreateBranch("b1", "master")
	assert.Nil(t, err)
	err = repo.CreateBranch("b2", "master")
	assert.Nil(t, err)
	err = origin.CreateBranch("b3", "master")
	assert.Nil(t, err)
	err = repo.Fetch()
	assert.Nil(t, err)
	branches, err := repo.Runner.RemoteBranches()
	assert.Nil(t, err)
	assert.Equal(t, []string{"origin/b3", "origin/master"}, branches)
}

func TestRunner_Remotes(t *testing.T) {
	runner := test.CreateRepo(t).Runner
	origin := test.CreateRepo(t)
	err := runner.AddRemote("origin", origin.WorkingDir())
	assert.Nil(t, err)
	remotes, err := runner.Remotes()
	assert.Nil(t, err)
	assert.Equal(t, []string{"origin"}, remotes)
}

func TestRunner_RemoveBranch(t *testing.T) {
	runner := test.CreateRepo(t).Runner
	err := runner.CreateBranch("b1", "master")
	assert.Nil(t, err)
	branches, err := runner.LocalBranches()
	assert.Nil(t, err)
	assert.Equal(t, []string{"b1", "master"}, branches)
	err = runner.RemoveBranch("b1")
	assert.Nil(t, err)
	branches, err = runner.LocalBranches()
	assert.Nil(t, err)
	assert.Equal(t, []string{"master"}, branches)
}

func TestRunner_RemoveRemote(t *testing.T) {
	runner := test.CreateRepo(t).Runner
	origin := test.CreateRepo(t)
	err := runner.AddRemote("origin", origin.WorkingDir())
	assert.Nil(t, err)
	err = runner.RemoveRemote("origin")
	assert.Nil(t, err)
	remotes, err := runner.Remotes()
	assert.Nil(t, err)
	assert.Len(t, remotes, 0)
}

func TestRunner_SetRemote(t *testing.T) {
	runner := test.CreateRepo(t).Runner
	remotes, err := runner.Remotes()
	assert.Nil(t, err)
	assert.Equal(t, []string{}, remotes)
	origin := test.CreateRepo(t)
	err = runner.AddRemote("origin", origin.WorkingDir())
	assert.Nil(t, err)
	remotes, err = runner.Remotes()
	assert.Nil(t, err)
	assert.Equal(t, []string{"origin"}, remotes)
}

func TestRunner_ShaForCommit(t *testing.T) {
	runner := test.CreateRepo(t).Runner
	err := runner.CreateCommit(git.Commit{Branch: "master", FileName: "foo", FileContent: "bar", Message: "commit"})
	assert.Nil(t, err)
	sha, err := runner.ShaForCommit("commit")
	assert.Nil(t, err)
	assert.Len(t, sha, 40)
}

func TestRunner_StageFile(t *testing.T) {
	runner := test.CreateRepo(t).Runner
	err := runner.CreateFile("f1.txt", "one")
	assert.Nil(t, err)
}

func TestRunner_Stash(t *testing.T) {
	runner := test.CreateRepo(t).Runner
	stashSize, err := runner.StashSize()
	assert.Nil(t, err)
	assert.Zero(t, stashSize)
	err = runner.CreateFile("f1.txt", "hello")
	assert.Nil(t, err)
	err = runner.Stash()
	assert.Nil(t, err)
	stashSize, err = runner.StashSize()
	assert.Nil(t, err)
	assert.Equal(t, 1, stashSize)
}

func TestRunner_UncommittedFiles(t *testing.T) {
	runner := test.CreateRepo(t).Runner
	err := runner.CreateFile("f1.txt", "one")
	assert.Nil(t, err)
	err = runner.CreateFile("f2.txt", "two")
	assert.Nil(t, err)
	files, err := runner.UncommittedFiles()
	assert.Nil(t, err)
	assert.Equal(t, []string{"f1.txt", "f2.txt"}, files)
}

// CreateTestGitTownRepo creates a GitRepo for use in tests, with a main branch and
// initial git town configuration
func CreateTestGitTownRepo(t *testing.T) test.Repo {
	repo := test.CreateRepo(t)
	err := repo.CreateBranch("main", "master")
	assert.Nil(t, err)
	err = repo.RunMany([][]string{
		{"git", "config", "git-town.main-branch-name", "main"},
		{"git", "config", "git-town.perennial-branch-names", ""},
	})
	assert.Nil(t, err)
	return repo
}
