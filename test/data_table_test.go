package test_test

import (
	"fmt"
	"testing"

	"github.com/git-town/git-town/test"
	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/stretchr/testify/assert"
)

func TestDataTable(t *testing.T) {
	r := test.DataTable{}
	r.AddRow("ALPHA", "BETA")
	r.AddRow("1", "2")
	r.AddRow("longer text", "even longer text")
	expected := `| ALPHA       | BETA             |
| 1           | 2                |
| longer text | even longer text |
`
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(expected, r.String(), false)
	if !(len(diffs) == 1 && diffs[0].Type == 0) {
		fmt.Println(dmp.DiffPrettyText(diffs))
		t.Fail()
	}
}

func TestDataTable_Remove(t *testing.T) {
	r := test.DataTable{}
	r.AddRow("local", "main, master, foo")
	r.AddRow("remote", "master, bar")
	r.RemoveText("master, ")
	expected := "| local  | main, foo |\n| remote | bar       |\n"
	assert.Equal(t, expected, r.String())
}
