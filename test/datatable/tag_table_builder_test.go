package datatable_test

import (
	"testing"

	"github.com/git-town/git-town/v13/test/datatable"
	"github.com/shoenig/test/must"
)

func TestTagTableBuilder(t *testing.T) {
	t.Parallel()
	builder := datatable.NewTagTableBuilder()
	builder.AddMany([]string{"tagB", "tagC"}, "local")
	builder.AddMany([]string{"tagA", "tagB"}, "origin")
	table := builder.Table()
	expected := `
| NAME | LOCATION      |
| tagA | origin        |
| tagB | local, origin |
| tagC | local         |
`
	must.EqOp(t, expected, "\n"+table.String())
}
