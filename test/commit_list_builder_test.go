package test

import (
	"testing"

	"github.com/Originate/git-town/test/gherkintools"
	"github.com/stretchr/testify/assert"
)

func TestCommitListBuilder(t *testing.T) {
	builder := NewCommitListBuilder()
	commit1 := gherkintools.Commit{SHA: "sha1", Branch: "branch1", Message: "commit1"}
	commit2 := gherkintools.Commit{SHA: "sha2", Branch: "branch2", Message: "commit2"}
	commit3 := gherkintools.Commit{SHA: "sha3", Branch: "branch3", Message: "commit3"}
	builder.Add(commit1, "local")
	builder.Add(commit1, "remote")
	builder.Add(commit2, "local")
	builder.Add(commit3, "remote")
	table := builder.Table([]string{"BRANCH", "LOCATION", "MESSAGE"})
	expected := `| BRANCH  | LOCATION      | MESSAGE |
| branch1 | local, remote | commit1 |
| branch2 | local         | commit2 |
| branch3 | remote        | commit3 |`
	assert.Equal(t, expected, table.String())
}
