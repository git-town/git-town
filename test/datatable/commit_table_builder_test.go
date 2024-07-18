package datatable_test

import (
	"testing"

	"github.com/git-town/git-town/v14/src/git/gitdomain"
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
	"github.com/git-town/git-town/v14/test/datatable"
	"github.com/git-town/git-town/v14/test/git"
	"github.com/shoenig/test/must"
)

func TestCommitTableBuilder(t *testing.T) {
	t.Parallel()
	builder := datatable.NewCommitTableBuilder()
	commit1 := git.Commit{SHA: Some(gitdomain.NewSHA("111111")), Branch: gitdomain.NewLocalBranchName("branch1"), Message: "commit1"}
	commit2 := git.Commit{SHA: Some(gitdomain.NewSHA("222222")), Branch: gitdomain.NewLocalBranchName("main"), Message: "commit2"}
	commit3 := git.Commit{SHA: Some(gitdomain.NewSHA("333333")), Branch: gitdomain.NewLocalBranchName("main"), Message: "commit3"}
	commit4 := git.Commit{SHA: Some(gitdomain.NewSHA("444444")), Branch: gitdomain.NewLocalBranchName("branch2"), Message: "commit4"}
	commit5 := git.Commit{SHA: Some(gitdomain.NewSHA("555555")), Branch: gitdomain.NewLocalBranchName("branch3"), Message: "commit5"}
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
