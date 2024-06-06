package datatable_test

import (
	"fmt"
	"testing"

	"github.com/git-town/git-town/v14/test/datatable"
	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/shoenig/test/must"
)

func TestDataTable(t *testing.T) {
	t.Parallel()

	t.Run("String serialization", func(t *testing.T) {
		t.Parallel()
		t.Run("normal table", func(t *testing.T) {
			t.Parallel()
			table := datatable.DataTable{}
			table.AddRow("ALPHA", "BETA")
			table.AddRow("1", "2")
			table.AddRow("longer text", "even longer text")
			expected := `| ALPHA       | BETA             |
| 1           | 2                |
| longer text | even longer text |
`
			dmp := diffmatchpatch.New()
			diffs := dmp.DiffMain(expected, table.String(), false)
			if !(len(diffs) == 1 && diffs[0].Type == 0) {
				fmt.Println(dmp.DiffPrettyText(diffs))
				t.Fail()
			}
		})
		t.Run("empty table", func(t *testing.T) {
			t.Parallel()
			table := datatable.DataTable{}
			want := ""
			have := table.String()
			must.EqOp(t, want, have)
		})
	})

	t.Run("RemoveText", func(t *testing.T) {
		t.Parallel()
		table := datatable.DataTable{}
		table.AddRow("local", "main, initial, foo")
		table.AddRow("origin", "initial, bar")
		table.RemoveText("initial, ")
		expected := "| local  | main, foo |\n| origin | bar       |\n"
		must.EqOp(t, expected, table.String())
	})

	t.Run("Sort", func(t *testing.T) {
		t.Parallel()
		table := datatable.DataTable{}
		table.AddRow("gamma", "3")
		table.AddRow("beta", "2")
		table.AddRow("alpha", "1")
		table.Sort()
		want := datatable.DataTable{Cells: [][]string{{"alpha", "1"}, {"beta", "2"}, {"gamma", "3"}}}
		diff, errCnt := table.EqualDataTable(want)
		if errCnt > 0 {
			t.Errorf("\nERROR! Found %d differences\n\n%s", errCnt, diff)
		}
	})

	t.Run("Strings", func(t *testing.T) {
		t.Parallel()
		table := datatable.DataTable{}
		table.AddRow("BRANCH", "TYPE", "COMMAND")
		table.AddRow("", "backend", "git version")
		table.AddRow("", "backend", "git config -lz --includes --global")
		table.AddRow("", "backend", "git config -lz --includes --local")
		table.AddRow("", "backend", "git rev-parse --show-toplevel")
		table.AddRow("", "backend", "git stash list")
		table.AddRow("", "backend", "git branch -vva")
		table.AddRow("", "backend", "git remote")
		table.AddRow("old", "frontend", "git fetch --prune --tags")
		table.AddRow("", "backend", "git branch -vva")
		table.AddRow("", "backend", "git status --long --ignore-submodules")
		table.AddRow("", "backend", "git rev-parse --verify --abbrev-ref @{-1}")
		table.AddRow("old", "frontend", "git merge --no-edit --ff main")
		table.AddRow("", "backend", "git diff main..old")
		table.AddRow("old", "frontend", "git checkout main")
		table.AddRow("main", "frontend", "git branch -D old")
		table.AddRow("", "backend", "git config git-town.perennial-branches")
		table.AddRow("", "backend", "git show-ref --quiet refs/heads/main")
		table.AddRow("", "backend", "git show-ref --quiet refs/heads/old")
		table.AddRow("", "backend", "git rev-parse --verify --abbrev-ref @{-1}")
		table.AddRow("", "backend", "git checkout main")
		table.AddRow("", "backend", "git checkout main")
		table.AddRow("", "backend", "git config -lz --includes --global")
		table.AddRow("", "backend", "git config -lz --includes --local")
		table.AddRow("", "backend", "git branch -vva")
		table.AddRow("", "backend", "git stash list")
		have := table.String()
		want := `
| BRANCH | TYPE     | COMMAND                                   |
|        | backend  | git version                               |
|        | backend  | git config -lz --includes --global                   |
|        | backend  | git config -lz --includes --local                    |
|        | backend  | git rev-parse --show-toplevel             |
|        | backend  | git stash list                            |
|        | backend  | git branch -vva                           |
|        | backend  | git remote                                |
| old    | frontend | git fetch --prune --tags                  |
|        | backend  | git branch -vva                           |
|        | backend  | git status --long --ignore-submodules     |
|        | backend  | git rev-parse --verify --abbrev-ref @{-1} |
| old    | frontend | git merge --no-edit --ff main             |
|        | backend  | git diff main..old                        |
| old    | frontend | git checkout main                         |
| main   | frontend | git branch -D old                         |
|        | backend  | git config git-town.perennial-branches    |
|        | backend  | git show-ref --quiet refs/heads/main      |
|        | backend  | git show-ref --quiet refs/heads/old       |
|        | backend  | git rev-parse --verify --abbrev-ref @{-1} |
|        | backend  | git checkout main                         |
|        | backend  | git checkout main                         |
|        | backend  | git config -lz --includes --global                   |
|        | backend  | git config -lz --includes --local                    |
|        | backend  | git branch -vva                           |
|        | backend  | git stash list                            |
`[1:]
		must.Eq(t, want, have)
	})
}
