package repo_test

import (
	"testing"

	"github.com/git-town/git-town/v8/test/git"
	"github.com/git-town/git-town/v8/test/repo"
	"github.com/stretchr/testify/assert"
)

func TestFileContentInCommit(t *testing.T) {
	t.Parallel()
	runtime := repo.Create(t)
	err := repo.CreateCommit(&runtime, git.Commit{
		Branch:      "initial",
		FileName:    "hello.txt",
		FileContent: "hello world",
		Message:     "commit",
	})
	assert.NoError(t, err)
	commits, err := repo.CommitsInBranch(&runtime, "initial", []string{})
	assert.NoError(t, err)
	assert.Len(t, commits, 1)
	content, err := repo.FileContentInCommit(&runtime, commits[0].SHA, "hello.txt")
	assert.NoError(t, err)
	assert.Equal(t, "hello world", content)
}
