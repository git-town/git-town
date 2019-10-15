package test

import (
	"testing"

	"github.com/Originate/git-town/test/gherkintools"
)

func TestCommitListBuilder(t *testing.T) {
	builder := NewCommitListBuilder()
	commit1 := gherkintools.Commit{SHA: "commit1", Branch: "branch1"}
	commit2 := gherkintools.Commit{SHA: "commit2", Branch: "branch2"}
	commit3 := gherkintools.Commit{SHA: "commit3", Branch: "branch3"}
	builder.Add(commit1, "local")
	builder.Add(commit1, "remote")
	builder.Add(commit2, "local")
	builder.Add(commit3, "remote")
	table := builder.Table([]string{"BRANCH", "LOCATION", "MESSAGE"})
	expected := `| BRANCH           | LOCATION      | MESSAGE                 |
| branch1      | local, remote | commit1 |
| branch2      | local, remote | commit2 |
}
