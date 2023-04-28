package datatable_test

import (
	"fmt"
	"testing"

	"github.com/git-town/git-town/v8/test/datatable"
	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/stretchr/testify/assert"
)

func TestDataTable(t *testing.T) {
	t.Parallel()
	t.Run("String serialization", func(t *testing.T) {
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

	t.Run("RemoveText", func(t *testing.T) {
		t.Parallel()
		table := datatable.DataTable{}
		table.AddRow("local", "main, initial, foo")
		table.AddRow("origin", "initial, bar")
		table.RemoveText("initial, ")
		expected := "| local  | main, foo |\n| origin | bar       |\n"
		assert.Equal(t, expected, table.String())
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
}
