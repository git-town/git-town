package datatable_test

import (
	"testing"

	"github.com/git-town/git-town/v20/internal/test/datatable"
	"github.com/git-town/git-town/v20/internal/test/testgit"
	"github.com/shoenig/test/must"
)

func TestCommitTableBuilder(t *testing.T) {
	t.Parallel()
	builder := datatable.NewCommitTableBuilder()
	commit1 := testgit.Commit{SHA: "111111", Branch: "branch1", Message: "commit1"}
	commit2 := testgit.Commit{SHA: "222222", Branch: "main", Message: "commit2"}
	commit3 := testgit.Commit{SHA: "333333", Branch: "main", Message: "commit3"}
	commit4 := testgit.Commit{SHA: "444444", Branch: "branch2", Message: "commit4"}
	commit5 := testgit.Commit{SHA: "555555", Branch: "branch3", Message: "commit5"}
	builder.Add(commit1, "local")
	builder.Add(commit1, "origin")
	builder.Add(commit1, "worktree")
	builder.Add(commit2, "local")
	builder.Add(commit3, "origin")
	builder.Add(commit4, "origin")
	builder.Add(commit4, "worktree")
	builder.Add(commit5, "worktree")
	table := builder.Table([]string{"BRANCH", "LOCATION", "MESSAGE"})
	expected := `
| BRANCH  | LOCATION                | MESSAGE |
| main    | local                   | commit2 |
|         | origin                  | commit3 |
| branch1 | local, origin, worktree | commit1 |
| branch2 | origin, worktree        | commit4 |
| branch3 | worktree                | commit5 |
`[1:]
	must.EqOp(t, expected, table.String())
}
