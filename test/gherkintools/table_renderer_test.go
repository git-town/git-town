package gherkintools

import (
	"fmt"
	"testing"

	"github.com/sergi/go-diff/diffmatchpatch"
)

func TestTableRenderer(t *testing.T) {
	r := TableRenderer{}
	r.AddLine("ALPHA", "BETA")
	r.AddLine("1", "2")
	r.AddLine("longer text", "even longer text")
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
