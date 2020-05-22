package git_test

import (
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/git-town/git-town/test"
	"github.com/stretchr/testify/assert"
)

func TestGitRepository_AddRemote(t *testing.T) {
	repo := CreateTestGitTownRepo(t)
	err := repo.AddRemote("foo", "bar")
	assert.Nil(t, err)
	remotes, err := repo.Remotes()
	assert.Nil(t, err)
	assert.Equal(t, []string{"foo"}, remotes)
}

func TestGitRepository_CheckoutBranch(t *testing.T) {
	repo := test.CreateTestRepo(t)
	err := repo.CreateBranch("branch1", "master")
	assert.Nil(t, err)
	err = repo.CheckoutBranch("branch1")
	assert.Nil(t, err)
	currentBranch, err := repo.CurrentBranch()
	assert.Nil(t, err)
	assert.Equal(t, "branch1", currentBranch)
	err = repo.CheckoutBranch("master")
	assert.Nil(t, err)
	currentBranch, err = repo.CurrentBranch()
	assert.Nil(t, err)
	assert.Equal(t, "master", currentBranch)
}

func TestGitRepository_Commits(t *testing.T) {
	repo := test.CreateTestRepo(t)
	err := repo.CreateCommit(test.Commit{
		Branch:      "master",
		FileName:    "file1",
		FileContent: "hello",
		Message:     "first commit",
	})
	assert.Nil(t, err)
	err = repo.CreateCommit(test.Commit{
		Branch:      "master",
		FileName:    "file2",
		FileContent: "hello again",
		Message:     "second commit",
	})
	assert.Nil(t, err)
	commits, err := repo.Commits([]string{"FILE NAME", "FILE CONTENT"})
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

func TestGitRepository_Configuration(t *testing.T) {
	repo := test.CreateTestRepo(t)
	config := repo.Config(false)
	assert.NotNil(t, config, "first path: new config")
	config = repo.Config(false)
	assert.NotNil(t, config, "second path: cached config")
}

func TestGitRepo_ConnectTrackingBranch(t *testing.T) {
	// replicating the situation this is used in,
	// connecting branches of repos with the same commits in them
	origin := test.CreateTestRepo(t)
	repoDir := filepath.Join(test.CreateTempDir(t), "repo") // need a non-existing directory
	err := test.CopyDirectory(origin.Dir, repoDir)
	assert.Nil(t, err)
	repo := test.NewGitRepository(repoDir, test.NewMockingShell(repoDir, repoDir, ""))
	err = repo.AddRemote("origin", origin.Dir)
	assert.Nil(t, err)
	err = repo.Fetch()
	assert.Nil(t, err)
	err = repo.ConnectTrackingBranch("master")
	assert.Nil(t, err)
	err = repo.PushBranch("master")
	assert.Nil(t, err)
}

func TestGitRepo_CreateBranch(t *testing.T) {
	repo := test.CreateTestRepo(t)
	err := repo.CreateBranch("branch1", "master")
	assert.Nil(t, err)
	currentBranch, err := repo.CurrentBranch()
	assert.Nil(t, err)
	assert.Equal(t, "master", currentBranch)
	branches, err := repo.LocalBranches()
	assert.Nil(t, err)
	assert.Equal(t, []string{"branch1", "master"}, branches)
}

func TestGitRepo_CreateChildFeatureBranch(t *testing.T) {
	repo := CreateTestGitTownRepo(t)
	err := repo.CreateFeatureBranch("f1")
	assert.Nil(t, err)
	err = repo.CreateChildFeatureBranch("f1a", "f1")
	assert.Nil(t, err)
	res, err := repo.Run("git", "town", "config")
	assert.Nil(t, err)
	has := strings.Contains(res.OutputSanitized(), "Branch Ancestry:\n  main\n    f1\n      f1a")
	assert.True(t, has)
}

func TestGitRepo_CreateCommit(t *testing.T) {
	repo := test.CreateTestRepo(t)
	err := repo.CreateCommit(test.Commit{
		Branch:      "master",
		FileName:    "hello.txt",
		FileContent: "hello world",
		Message:     "test commit",
	})
	assert.Nil(t, err)
	commits, err := repo.Commits([]string{"FILE NAME", "FILE CONTENT"})
	assert.Nil(t, err)
	assert.Len(t, commits, 1)
	assert.Equal(t, "hello.txt", commits[0].FileName)
	assert.Equal(t, "hello world", commits[0].FileContent)
	assert.Equal(t, "test commit", commits[0].Message)
	assert.Equal(t, "master", commits[0].Branch)
}

func TestGitRepo_CreateCommit_Author(t *testing.T) {
	repo := test.CreateTestRepo(t)
	err := repo.CreateCommit(test.Commit{
		Branch:      "master",
		FileName:    "hello.txt",
		FileContent: "hello world",
		Message:     "test commit",
		Author:      "developer <developer@example.com>",
	})
	assert.Nil(t, err)
	commits, err := repo.Commits([]string{"FILE NAME", "FILE CONTENT"})
	assert.Nil(t, err)
	assert.Len(t, commits, 1)
	assert.Equal(t, "hello.txt", commits[0].FileName)
	assert.Equal(t, "hello world", commits[0].FileContent)
	assert.Equal(t, "test commit", commits[0].Message)
	assert.Equal(t, "master", commits[0].Branch)
	assert.Equal(t, "developer <developer@example.com>", commits[0].Author)
}

func TestGitRepo_CreateFeatureBranch(t *testing.T) {
	repo := CreateTestGitTownRepo(t)
	err := repo.CreateFeatureBranch("f1")
	assert.Nil(t, err)
	assert.True(t, repo.Config(true).IsFeatureBranch("f1"))
	assert.Equal(t, []string{"main"}, repo.Config(true).GetAncestorBranches("f1"))
}

func TestGitRepo_CreateFeatureBranchNoParent(t *testing.T) {
	repo := CreateTestGitTownRepo(t)
	err := repo.CreateFeatureBranchNoParent("f1")
	assert.Nil(t, err)
	assert.True(t, repo.Config(true).IsFeatureBranch("f1"))
	assert.Equal(t, []string(nil), repo.Config(true).GetAncestorBranches("f1"))
}

func TestGitRepo_CreateFile(t *testing.T) {
	repo := test.CreateTestRepo(t)
	err := repo.CreateFile("filename", "content")
	assert.Nil(t, err, "cannot create file in repo")
	content, err := ioutil.ReadFile(filepath.Join(repo.Dir, "filename"))
	assert.Nil(t, err, "cannot read file")
	assert.Equal(t, "content", string(content))
}

func TestGitRepo_CreateFile_InSubFolder(t *testing.T) {
	repo := test.CreateTestRepo(t)
	err := repo.CreateFile("folder/filename", "content")
	assert.Nil(t, err, "cannot create file in repo")
	content, err := ioutil.ReadFile(filepath.Join(repo.Dir, "folder/filename"))
	assert.Nil(t, err, "cannot read file")
	assert.Equal(t, "content", string(content))
}

func TestGitRepo_CreatePerennialBranches(t *testing.T) {
	repo := CreateTestGitTownRepo(t)
	err := repo.CreatePerennialBranches("p1", "p2")
	assert.Nil(t, err)
	branches, err := repo.LocalBranches()
	assert.Nil(t, err)
	assert.Equal(t, []string{"main", "master", "p1", "p2"}, branches)
	config := repo.Config(true)
	assert.True(t, config.IsPerennialBranch("p1"))
	assert.True(t, config.IsPerennialBranch("p2"))
}

func TestGitRepo_CurrentBranch(t *testing.T) {
	repo := test.CreateTestRepo(t)
	err := repo.CheckoutBranch("master")
	assert.Nil(t, err)
	err = repo.CreateBranch("b1", "master")
	assert.Nil(t, err)
	err = repo.CheckoutBranch("b1")
	assert.Nil(t, err)
	branch, err := repo.CurrentBranch()
	assert.Nil(t, err)
	assert.Equal(t, "b1", branch)
	err = repo.CheckoutBranch("master")
	assert.Nil(t, err)
	branch, err = repo.CurrentBranch()
	assert.Nil(t, err)
	assert.Equal(t, "master", branch)
}

func TestGitRepo_Fetch(t *testing.T) {
	repo := test.CreateTestRepo(t)
	origin := test.CreateTestRepo(t)
	err := repo.AddRemote("origin", origin.Dir)
	assert.Nil(t, err)
	err = repo.Fetch()
	assert.Nil(t, err)
}

func TestGitRepo_FilesInCommit(t *testing.T) {
	repo := test.CreateTestRepo(t)
	err := repo.CreateFile("f1.txt", "one")
	assert.Nil(t, err)
	err = repo.CreateFile("f2.txt", "two")
	assert.Nil(t, err)
	err = repo.StageFiles("f1.txt", "f2.txt")
	assert.Nil(t, err)
	err = repo.CommitStagedChanges(true)
	assert.Nil(t, err)
	commits, err := repo.Commits([]string{})
	assert.Nil(t, err)
	assert.Len(t, commits, 1)
	fileNames, err := repo.FilesInCommit(commits[0].SHA)
	assert.Nil(t, err)
	assert.Equal(t, []string{"f1.txt", "f2.txt"}, fileNames)
}

func TestGitRepo_HasBranchesOutOfSync_synced(t *testing.T) {
	repo1 := test.CreateTestRepo(t)
	dir2 := test.CreateTempDir(t)
	repo2, err := test.CloneGitRepo(repo1.Dir, dir2, repo1.Dir, repo1.Dir)
	assert.Nil(t, err)
	err = repo2.CreateBranch("branch1", "master")
	assert.Nil(t, err)
	err = repo2.CheckoutBranch("branch1")
	assert.Nil(t, err)
	err = repo2.CreateFile("file1", "content")
	assert.Nil(t, err)
	err = repo2.StageFiles("file1")
	assert.Nil(t, err)
	err = repo2.CommitStagedChanges(true)
	assert.Nil(t, err)
	err = repo2.PushBranch("master")
	assert.Nil(t, err)
	have, err := repo2.HasBranchesOutOfSync()
	assert.Nil(t, err)
	assert.False(t, have)
}

func TestGitRepo_HasBranchesOutOfSync_branchAhead(t *testing.T) {
	repo1 := test.CreateTestRepo(t)
	dir2 := test.CreateTempDir(t)
	repo2, err := test.CloneGitRepo(repo1.Dir, dir2, repo1.Dir, repo1.Dir)
	assert.Nil(t, err)
	err = repo2.CreateBranch("branch1", "master")
	assert.Nil(t, err)
	err = repo2.PushBranch("branch1")
	assert.Nil(t, err)
	err = repo2.CreateFile("file1", "content")
	assert.Nil(t, err)
	err = repo2.StageFiles("file1")
	assert.Nil(t, err)
	err = repo2.CommitStagedChanges(true)
	assert.Nil(t, err)
	have, err := repo2.HasBranchesOutOfSync()
	assert.Nil(t, err)
	assert.True(t, have)
}

func TestGitRepo_HasBranchesOutOfSync_branchBehind(t *testing.T) {
	repo1 := test.CreateTestRepo(t)
	dir2 := test.CreateTempDir(t)
	repo2, err := test.CloneGitRepo(repo1.Dir, dir2, repo1.Dir, repo1.Dir)
	assert.Nil(t, err)
	err = repo2.CreateBranch("branch1", "master")
	assert.Nil(t, err)
	err = repo2.PushBranch("branch1")
	assert.Nil(t, err)
	err = repo1.CreateFile("file1", "content")
	assert.Nil(t, err)
	err = repo1.StageFiles("file1")
	assert.Nil(t, err)
	err = repo1.CommitStagedChanges(true)
	assert.Nil(t, err)
	err = repo2.Fetch()
	assert.Nil(t, err)
	have, err := repo2.HasBranchesOutOfSync()
	assert.Nil(t, err)
	assert.True(t, have)
}

func TestGitRepo_HasGitTownConfigNow(t *testing.T) {
	repo := test.CreateTestRepo(t)
	res, err := repo.HasGitTownConfigNow()
	assert.Nil(t, err)
	assert.False(t, res)
	err = repo.CreateBranch("main", "master")
	assert.Nil(t, err)
	err = repo.CreateFeatureBranch("foo")
	assert.Nil(t, err)
	res, err = repo.HasGitTownConfigNow()
	assert.Nil(t, err)
	assert.True(t, res)
}

func TestGitRepo_HasFile(t *testing.T) {
	repo := test.CreateTestRepo(t)
	err := repo.CreateFile("f1.txt", "one")
	assert.Nil(t, err)
	has, err := repo.HasFile("f1.txt", "one")
	assert.Nil(t, err)
	assert.True(t, has)
	_, err = repo.HasFile("f1.txt", "zonk")
	assert.Error(t, err)
	_, err = repo.HasFile("zonk.txt", "one")
	assert.Error(t, err)
}
func TestGitRepo_HasLocalBranch(t *testing.T) {
	origin := test.CreateTestRepo(t)
	repoDir := test.CreateTempDir(t)
	repo, err := test.CloneGitRepo(origin.Dir, repoDir, repoDir, repoDir)
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

func TestGitRepo_HasRebaseInProgress(t *testing.T) {
	repo := test.CreateTestRepo(t)
	has, err := repo.HasRebaseInProgress()
	assert.Nil(t, err)
	assert.False(t, has)
}

func TestGitRepo_HasRemote(t *testing.T) {
	origin := test.CreateTestRepo(t)
	repoDir := test.CreateTempDir(t)
	repo, err := test.CloneGitRepo(origin.Dir, repoDir, repoDir, repoDir)
	assert.Nil(t, err)
	has, err := repo.HasRemote("origin")
	assert.Nil(t, err)
	assert.True(t, has)
	has, err = repo.HasRemote("zonk")
	assert.Nil(t, err)
	assert.False(t, has)
}

func TestGitRepo_LastActiveDir(t *testing.T) {
	repo := test.CreateTestRepo(t)
	dir, err := repo.LastActiveDir()
	assert.Nil(t, err)
	assert.Equal(t, repo.Dir, dir)
}

func TestGitRepo_LocalBranches(t *testing.T) {
	origin := test.CreateTestRepo(t)
	repoDir := test.CreateTempDir(t)
	repo, err := test.CloneGitRepo(origin.Dir, repoDir, repoDir, repoDir)
	assert.Nil(t, err)
	err = repo.CreateBranch("b1", "master")
	assert.Nil(t, err)
	err = repo.CreateBranch("b2", "master")
	assert.Nil(t, err)
	err = origin.CreateBranch("b3", "master")
	assert.Nil(t, err)
	err = repo.Fetch()
	assert.Nil(t, err)
	branches, err := repo.LocalBranches()
	assert.Nil(t, err)
	assert.Equal(t, []string{"b1", "b2", "master"}, branches)
}

func TestGitRepo_LocalAndRemoteBranches(t *testing.T) {
	origin := test.CreateTestRepo(t)
	repoDir := test.CreateTempDir(t)
	repo, err := test.CloneGitRepo(origin.Dir, repoDir, repoDir, repoDir)
	assert.Nil(t, err)
	err = repo.CreateBranch("b1", "master")
	assert.Nil(t, err)
	err = repo.CreateBranch("b2", "master")
	assert.Nil(t, err)
	err = origin.CreateBranch("b3", "master")
	assert.Nil(t, err)
	err = repo.Fetch()
	assert.Nil(t, err)
	branches, err := repo.LocalAndRemoteBranches()
	assert.Nil(t, err)
	assert.Equal(t, []string{"b1", "b2", "b3", "master"}, branches)
}

func TestGitRepo_PushBranch(t *testing.T) {
	repo := test.CreateTestRepo(t)
	origin := test.CreateTestRepo(t)
	err := repo.AddRemote("origin", origin.Dir)
	assert.Nil(t, err)
	err = repo.CreateBranch("b1", "master")
	assert.Nil(t, err)
	err = repo.PushBranch("b1")
	assert.Nil(t, err)
	branches, err := origin.LocalBranches()
	assert.Nil(t, err)
	assert.Equal(t, []string{"b1", "master"}, branches)
}

func TestGitRepo_Remotes(t *testing.T) {
	repo := test.CreateTestRepo(t)
	origin := test.CreateTestRepo(t)
	err := repo.AddRemote("origin", origin.Dir)
	assert.Nil(t, err)
	remotes, err := repo.Remotes()
	assert.Nil(t, err)
	assert.Equal(t, []string{"origin"}, remotes)
}

func TestGitRepo_RemoveBranch(t *testing.T) {
	repo := test.CreateTestRepo(t)
	err := repo.CreateBranch("b1", "master")
	assert.Nil(t, err)
	branches, err := repo.LocalBranches()
	assert.Nil(t, err)
	assert.Equal(t, []string{"b1", "master"}, branches)
	err = repo.RemoveBranch("b1")
	assert.Nil(t, err)
	branches, err = repo.LocalBranches()
	assert.Nil(t, err)
	assert.Equal(t, []string{"master"}, branches)
}

func TestGitRepo_RemoveRemote(t *testing.T) {
	repo := test.CreateTestRepo(t)
	origin := test.CreateTestRepo(t)
	err := repo.AddRemote("origin", origin.Dir)
	assert.Nil(t, err)
	err = repo.RemoveRemote("origin")
	assert.Nil(t, err)
	remotes, err := repo.Remotes()
	assert.Nil(t, err)
	assert.Len(t, remotes, 0)
}

func TestGitRepo_SetOffline(t *testing.T) {
	repo := CreateTestGitTownRepo(t)
	err := repo.SetOffline(true)
	assert.Nil(t, err)
	offline, err := repo.IsOffline()
	assert.Nil(t, err)
	assert.True(t, offline)
	err = repo.SetOffline(false)
	assert.Nil(t, err)
	offline, err = repo.IsOffline()
	assert.Nil(t, err)
	assert.False(t, offline)
}

func TestGitRepo_SetRemote(t *testing.T) {
	repo := test.CreateTestRepo(t)
	remotes, err := repo.Remotes()
	assert.Nil(t, err)
	assert.Equal(t, []string{}, remotes)
	origin := test.CreateTestRepo(t)
	err = repo.AddRemote("origin", origin.Dir)
	assert.Nil(t, err)
	remotes, err = repo.Remotes()
	assert.Nil(t, err)
	assert.Equal(t, []string{"origin"}, remotes)
}

func TestGitRepo_ShaForCommit(t *testing.T) {
	repo := test.CreateTestRepo(t)
	err := repo.CreateCommit(test.Commit{Branch: "master", FileName: "foo", FileContent: "bar", Message: "commit"})
	assert.Nil(t, err)
	sha, err := repo.ShaForCommit("commit")
	assert.Nil(t, err)
	assert.Len(t, sha, 40)
}

func TestGitRepo_StageFile(t *testing.T) {
	repo := test.CreateTestRepo(t)
	err := repo.CreateFile("f1.txt", "one")
	assert.Nil(t, err)
}

func TestGitRepo_Stash(t *testing.T) {
	repo := test.CreateTestRepo(t)
	stashSize, err := repo.StashSize()
	assert.Nil(t, err)
	assert.Zero(t, stashSize)
	err = repo.CreateFile("f1.txt", "hello")
	assert.Nil(t, err)
	err = repo.Stash()
	assert.Nil(t, err)
	stashSize, err = repo.StashSize()
	assert.Nil(t, err)
	assert.Equal(t, 1, stashSize)
}

func TestGitRepo_UncommittedFiles(t *testing.T) {
	repo := test.CreateTestRepo(t)
	err := repo.CreateFile("f1.txt", "one")
	assert.Nil(t, err)
	err = repo.CreateFile("f2.txt", "two")
	assert.Nil(t, err)
	files, err := repo.UncommittedFiles()
	assert.Nil(t, err)
	assert.Equal(t, []string{".gitconfig", "f1.txt", "f2.txt"}, files)
}

// CreateTestGitTownRepo creates a GitRepo for use in tests, with a main branch and
// initial git town configuration
func CreateTestGitTownRepo(t *testing.T) test.GitRepo {
	repo := test.CreateTestRepo(t)
	err := repo.CreateBranch("main", "master")
	assert.Nil(t, err)
	err = repo.RunMany([][]string{
		{"git", "config", "git-town.main-branch-name", "main"},
		{"git", "config", "git-town.perennial-branch-names", ""},
	})
	assert.Nil(t, err)
	return repo
}
