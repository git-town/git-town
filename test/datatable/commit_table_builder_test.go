package datatable_test

import (
	"testing"

	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/test/datatable"
	"github.com/git-town/git-town/v9/test/git"
	"github.com/stretchr/testify/assert"
)

func TestCommitTableBuilder(t *testing.T) {
	t.Parallel()
	builder := datatable.NewCommitTableBuilder()
	commit1 := git.Commit{SHA: domain.NewSHA("111111"), Branch: domain.NewLocalBranchName("branch1"), Message: "commit1"}
	commit2 := git.Commit{SHA: domain.NewSHA("222222"), Branch: domain.NewLocalBranchName("main"), Message: "commit2"}
	commit3 := git.Commit{SHA: domain.NewSHA("333333"), Branch: domain.NewLocalBranchName("main"), Message: "commit3"}
	commit4 := git.Commit{SHA: domain.NewSHA("444444"), Branch: domain.NewLocalBranchName("branch3"), Message: "commit4"}
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
