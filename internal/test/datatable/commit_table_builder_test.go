package datatable_test

import (
	"testing"

	"github.com/git-town/git-town/v24/internal/config/configdomain"
	"github.com/git-town/git-town/v24/internal/test/datatable"
	"github.com/git-town/git-town/v24/internal/test/testgit"
	"github.com/shoenig/test/must"
)

func TestCommitTableBuilder(t *testing.T) {
	t.Parallel()
	builder := datatable.NewCommitTableBuilder()
	commit1 := testgit.Commit{Branch: "branch1", Message: "commit1", SHA: "111111"}
	commit2 := testgit.Commit{Branch: "main", Message: "commit2", SHA: "222222"}
	commit3 := testgit.Commit{Branch: "main", Message: "commit3", SHA: "333333"}
	commit4 := testgit.Commit{Branch: "branch2", Message: "commit4", SHA: "444444"}
	commit5 := testgit.Commit{Branch: "branch3", Message: "commit5", SHA: "555555"}
	builder.Add(commit1, "local")
	builder.Add(commit1, "origin")
	builder.Add(commit1, "worktree")
	builder.Add(commit2, "local")
	builder.Add(commit3, "origin")
	builder.Add(commit4, "origin")
	builder.Add(commit4, "worktree")
	builder.Add(commit5, "worktree")
	lineage := configdomain.NewLineageWith(configdomain.LineageData{
		"branch1": "main",
		"branch2": "branch1",
		"branch3": "branch2",
	})
	table := builder.Table([]string{"BRANCH", "LOCATION", "MESSAGE"}, lineage, configdomain.OrderAsc)
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
