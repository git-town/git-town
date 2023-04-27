package commands_test

import (
	"testing"

	"github.com/git-town/git-town/v8/test/commands"
	"github.com/git-town/git-town/v8/test/git"
	"github.com/git-town/git-town/v8/test/runtime"
	"github.com/stretchr/testify/assert"
)

func TestShaForCommit(t *testing.T) {
	t.Parallel()
	repo := runtime.Create(t)
	err := commands.CreateCommit(&repo, git.Commit{Branch: "initial", FileName: "foo", FileContent: "bar", Message: "commit"})
	assert.NoError(t, err)
	sha, err := commands.ShaForCommit(&repo, "commit")
	assert.NoError(t, err)
	assert.Len(t, sha, 40)
}
