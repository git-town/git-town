package test

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCloneGitRepository(t *testing.T) {
	rootDir := createTempDir(t)
	originPath := filepath.Join(rootDir, "origin")
	_, err := InitGitRepository(originPath, rootDir)
	assert.Nil(t, err, "cannot initialze origin Git repository")
	clonedPath := filepath.Join(rootDir, "cloned")
	_, err = CloneGitRepository(originPath, clonedPath, rootDir)
	assert.Nil(t, err, "cannot clone repo")
	assertIsNormalGitRepo(t, clonedPath)
}

func TestInitGitRepository(t *testing.T) {
	dir := createTempDir(t)
	repo, err := InitGitRepository(dir, dir)
	assert.Nil(t, err, "cannot initialize normal GitRepository")
	assertIsNormalGitRepo(t, repo.Dir)
	// ensure the Git repo works, i.e. we can commit into it
	repo.CreateFile("test.txt", "hello")
	repo.StageFile("test.txt")
	repo.CommitStagedChanges()
}

func TestNewGitRepository(t *testing.T) {
	dir := createTempDir(t)
	_ = NewGitRepository(dir, dir)
}

func TestGitRepository_Branches(t *testing.T) {
	repo := createTestGitTownRepo(t)
	assert.Nil(t, repo.CreateFeatureBranch("branch3", false), "cannot create branch3")
	assert.Nil(t, repo.CreateFeatureBranch("branch2", false), "cannot create branch2")
	assert.Nil(t, repo.CreateFeatureBranch("branch1", false), "cannot create branch1")
	branches, err := repo.Branches()
	assert.Nil(t, err)
	assert.Equal(t, []string{"branch1", "branch2", "branch3", "master"}, branches)
}

func TestGitRepository_CheckoutBranch(t *testing.T) {
	repo := createTestRepo(t)
	repo.CreateBranch("branch1")
	repo.CheckoutBranch("branch1")
	currentBranch, err := repo.CurrentBranch()
	assert.Nil(t, err)
	assert.Equal(t, "branch1", currentBranch)
	repo.CheckoutBranch("master")
	currentBranch, err = repo.CurrentBranch()
	assert.Nil(t, err)
	assert.Equal(t, "master", currentBranch)
}

func TestGitRepository_Commits(t *testing.T) {
	repo := createTestRepo(t)
	repo.CreateCommit(Commit{
		Branch:      "master",
		FileName:    "file1",
		FileContent: "hello",
		Message:     "first commit",
	})
	repo.CreateCommit(Commit{
		Branch:      "master",
		FileName:    "file2",
		FileContent: "hello again",
		Message:     "second commit",
	})
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
	repo := createTestRepo(t)
	config := repo.Configuration(false)
	assert.NotNil(t, config, "first path: new config")
	config = repo.Configuration(false)
	assert.NotNil(t, config, "second path: cached config")
}

func TestGitRepo_ConnectTrackingBranch(t *testing.T) {
	repo := createTestRepo(t)
	origin := createTestRepo(t)
	err := repo.SetRemote(origin.homeDir)
	assert.Nil(t, err)
	err = repo.Fetch()
	assert.Nil(t, err)
	err = repo.ConnectTrackingBranch("master")
	assert.Nil(t, err)
	err = repo.PushBranch("master")
	assert.Nil(t, err)
}

func TestGitRepo_CreateBranch(t *testing.T) {
	repo := createTestRepo(t)
	err := repo.CreateBranch("branch1")
	assert.Nil(t, err)
	branches, err := repo.Branches()
	assert.Nil(t, err)
	assert.Equal(t, []string{"branch1", "master"}, branches)
}

func TestGitRepo_CreateChildFeatureBranch(t *testing.T) {
	repo := createTestGitTownRepo(t)
	err := repo.CreatePerennialBranches("main")
	assert.Nil(t, err)
	err = repo.CreateFeatureBranch("f1", false)
	assert.Nil(t, err)
	err = repo.CreateChildFeatureBranch("f1a", "f1")
	assert.Nil(t, err)
	branches, err := repo.Branches()
	assert.Nil(t, err)
	assert.Equal(t, []string{"f1", "f1a", "main", "master"}, branches)
}

func TestGitRepository_CreateFile(t *testing.T) {
	repo := createTestRepo(t)
	err := repo.CreateFile("filename", "content")
	assert.Nil(t, err, "cannot create file in repo")
	content, err := ioutil.ReadFile(filepath.Join(repo.Dir, "filename"))
	assert.Nil(t, err, "cannot read file")
	assert.Equal(t, "content", string(content))
}

func TestGitRepository_CreateFeatureBranch(t *testing.T) {
	repo := createTestGitTownRepo(t)
	err := repo.CreatePerennialBranches("main")
	assert.Nil(t, err)
	err = repo.CreateFeatureBranch("f1", false)
	assert.Nil(t, err)
	assert.True(t, repo.Configuration(true).IsFeatureBranch("f1"))
}

func TestGitRepository_CreateFile_InSubFolder(t *testing.T) {
	repo := createTestRepo(t)
	err := repo.CreateFile("folder/filename", "content")
	assert.Nil(t, err, "cannot create file in repo")
	content, err := ioutil.ReadFile(filepath.Join(repo.Dir, "folder/filename"))
	assert.Nil(t, err, "cannot read file")
	assert.Equal(t, "content", string(content))
}

// HELPERS

// createTestGitRepo creates a fully initialized Git repo including a master branch.
func createTestRepo(t *testing.T) GitRepository {
	dir := createTempDir(t)
	repo, err := InitGitRepository(dir, dir)
	assert.Nil(t, err, "cannot initialize Git repow")
	err = repo.RunMany([][]string{
		{"git", "commit", "--allow-empty", "-m", "initial commit"},
	})
	assert.Nil(t, err, "cannot create initial commit: %s")
	return repo
}

func createTestGitTownRepo(t *testing.T) GitRepository {
	repo := createTestRepo(t)
	err := repo.RunMany([][]string{
		{"git", "config", "git-town.main-branch-name", "master"},
		{"git", "config", "git-town.perennial-branch-names", ""},
	})
	assert.Nil(t, err)
	return repo
}
