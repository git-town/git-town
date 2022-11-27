package test_test

import (
	"fmt"
	"testing"

	"github.com/git-town/git-town/v7/test"
	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/stretchr/testify/assert"
)

func TestDataTable(t *testing.T) {
	t.Run("String serialization", func(t *testing.T) {
		t.Parallel()
		table := test.DataTable{}
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

	t.Run("RemoveText", func(t *testing.T) {
		t.Parallel()
		table := test.DataTable{}
		table.AddRow("local", "main, master, foo")
		table.AddRow("origin", "master, bar")
		table.RemoveText("master, ")
		expected := "| local  | main, foo |\n| origin | bar       |\n"
		assert.Equal(t, expected, table.String())
	})

	t.Run("Sort", func(t *testing.T) {
		t.Parallel()
		table := test.DataTable{}
		table.AddRow("gamma", "3")
		table.AddRow("beta", "2")
		table.AddRow("alpha", "1")
		table.Sort()
		want := test.DataTable{Cells: [][]string{{"alpha", "1"}, {"beta", "2"}, {"gamma", "3"}}}
		diff, errCnt := table.EqualDataTable(want)
		if errCnt > 0 {
			t.Errorf("\nERROR! Found %d differences\n\n%s", errCnt, diff)
		}
	})
}
