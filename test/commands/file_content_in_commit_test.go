package commands_test

import (
	"testing"

	"github.com/git-town/git-town/v8/test/commands"
	"github.com/git-town/git-town/v8/test/git"
	"github.com/stretchr/testify/assert"
)

func TestFileContentInCommit(t *testing.T) {
	t.Parallel()
	runtime := commands.Create(t)
	err := commands.CreateCommit(&runtime, git.Commit{
		Branch:      "initial",
		FileName:    "hello.txt",
		FileContent: "hello world",
		Message:     "commit",
	})
	assert.NoError(t, err)
	commits, err := commands.CommitsInBranch(&runtime, "initial", []string{})
	assert.NoError(t, err)
	assert.Len(t, commits, 1)
	content, err := commands.FileContentInCommit(&runtime, commits[0].SHA, "hello.txt")
	assert.NoError(t, err)
	assert.Equal(t, "hello world", content)
}
