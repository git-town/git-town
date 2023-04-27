package repo_test

import (
	"testing"

	"github.com/git-town/git-town/v8/test/git"
	"github.com/git-town/git-town/v8/test/repo"
	"github.com/stretchr/testify/assert"
)

func TestCommits(t *testing.T) {
	t.Parallel()
	dev := repo.Create(t)
	err := repo.CreateCommit(&dev, git.Commit{
		Branch:      "initial",
		FileName:    "file1",
		FileContent: "hello",
		Message:     "first commit",
	})
	assert.NoError(t, err)
	err = repo.CreateCommit(&dev, git.Commit{
		Branch:      "initial",
		FileName:    "file2",
		FileContent: "hello again",
		Message:     "second commit",
	})
	assert.NoError(t, err)
	commits, err := repo.Commits(&dev, []string{"FILE NAME", "FILE CONTENT"}, "initial")
	assert.NoError(t, err)
	assert.Len(t, commits, 2)
	assert.Equal(t, "initial", commits[0].Branch)
	assert.Equal(t, "file1", commits[0].FileName)
	assert.Equal(t, "hello", commits[0].FileContent)
	assert.Equal(t, "first commit", commits[0].Message)
	assert.Equal(t, "initial", commits[1].Branch)
	assert.Equal(t, "file2", commits[1].FileName)
	assert.Equal(t, "hello again", commits[1].FileContent)
	assert.Equal(t, "second commit", commits[1].Message)
}
