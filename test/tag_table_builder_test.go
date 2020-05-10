package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTagTableBuilder(t *testing.T) {
	builder := NewTagTableBuilder()
	builder.AddMany([]string{"tagB", "tagC"}, "local")
	builder.AddMany([]string{"tagA", "tagB"}, "remote")
	table := builder.Table()
	expected := `
| NAME | LOCATION      |
| tagA | remote        |
| tagB | local, remote |
| tagC | local         |
`
	assert.Equal(t, expected, "\n"+table.String())
}
