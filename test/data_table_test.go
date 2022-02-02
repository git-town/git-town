package test_test

import (
	"fmt"
	"testing"

	"github.com/git-town/git-town/v7/test"
	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/stretchr/testify/assert"
)

func TestDataTable(t *testing.T) {
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
}

func TestDataTable_Remove(t *testing.T) {
	t.Parallel()
	table := test.DataTable{}
	table.AddRow("local", "main, master, foo")
	table.AddRow("remote", "master, bar")
	table.RemoveText("master, ")
	expected := "| local  | main, foo |\n| remote | bar       |\n"
	assert.Equal(t, expected, table.String())
}
