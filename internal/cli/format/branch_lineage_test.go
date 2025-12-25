package format_test

import (
	"testing"

	"github.com/git-town/git-town/v22/internal/cli/format"
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/shoenig/test/must"
)

func TestBranchLineage(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		lineage := configdomain.NewLineage()
		have := format.BranchLineage(lineage, configdomain.OrderAsc)
		want := ""
		must.EqOp(t, want, have)
	})

	t.Run("multiple roots", func(t *testing.T) {
		t.Parallel()
		lineage := configdomain.NewLineageWith(configdomain.LineageData{
			"branch-1":  "main",
			"branch-1A": "branch-1",
			"branch-1B": "branch-1",
			"branch-2":  "main",
			"hotfix":    "qa",
		})
		have := format.BranchLineage(lineage, configdomain.OrderAsc)
		want := `
main
  branch-1
    branch-1A
    branch-1B
  branch-2

qa
  hotfix`[1:]
		must.EqOp(t, want, have)
	})

	t.Run("normal", func(t *testing.T) {
		t.Parallel()
		lineage := configdomain.NewLineageWith(configdomain.LineageData{
			"branch-1":  "main",
			"branch-1A": "branch-1",
			"branch-1B": "branch-1",
			"branch-2":  "main",
		})
		have := format.BranchLineage(lineage, configdomain.OrderAsc)
		want := `
main
  branch-1
    branch-1A
    branch-1B
  branch-2`[1:]
		must.EqOp(t, want, have)
	})
}
