package test

import (
	"fmt"
	"testing"

	"github.com/sergi/go-diff/diffmatchpatch"
)

func TestMortadella(t *testing.T) {
	r := Mortadella{}
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
