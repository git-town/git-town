package repo_test

import (
	"testing"

	"github.com/git-town/git-town/v8/test/git"
	"github.com/git-town/git-town/v8/test/repo"
	"github.com/stretchr/testify/assert"
)

func TestShaForCommit(t *testing.T) {
	t.Parallel()
	dev := repo.Create(t)
	err := repo.CreateCommit(&dev, git.Commit{Branch: "initial", FileName: "foo", FileContent: "bar", Message: "commit"})
	assert.NoError(t, err)
	sha, err := repo.ShaForCommit(&dev, "commit")
	assert.NoError(t, err)
	assert.Len(t, sha, 40)
}
