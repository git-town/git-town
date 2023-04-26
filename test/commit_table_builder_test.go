package test_test

import (
	"testing"

	"github.com/git-town/git-town/v8/test"
	"github.com/git-town/git-town/v8/test/git"
	"github.com/stretchr/testify/assert"
)

func TestCommitTableBuilder(t *testing.T) {
	t.Parallel()
	builder := test.NewCommitTableBuilder()
	commit1 := git.Commit{SHA: "sha1", Branch: "branch1", Message: "commit1"}
	commit2 := git.Commit{SHA: "sha2", Branch: "main", Message: "commit2"}
	commit3 := git.Commit{SHA: "sha3", Branch: "main", Message: "commit3"}
	commit4 := git.Commit{SHA: "sha4", Branch: "branch3", Message: "commit4"}
	builder.Add(commit1, "local")
	builder.Add(commit1, "origin")
	builder.Add(commit2, "local")
	builder.Add(commit3, "origin")
	builder.Add(commit4, "origin")
	table := builder.Table([]string{"BRANCH", "LOCATION", "MESSAGE"})
	expected := `| BRANCH  | LOCATION      | MESSAGE |
| main    | local         | commit2 |
|         | origin        | commit3 |
| branch1 | local, origin | commit1 |
| branch3 | origin        | commit4 |
`
	assert.Equal(t, expected, table.String())
}
