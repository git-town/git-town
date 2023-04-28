package datatable_test

import (
	"testing"

	"github.com/git-town/git-town/v8/test/datatable"
	"github.com/stretchr/testify/assert"
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
	assert.Equal(t, expected, "\n"+table.String())
}
