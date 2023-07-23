package output

import (
	"github.com/git-town/git-town/v9/test/datatable"
)

// RenderExecutedGitCommands provides the textual Gherkin table representation of the given executed Git commands.
// The DataTable table matches the structure of the given Gherkin table.
func RenderExecutedDebugCommands(commands []string) datatable.DataTable {
	result := datatable.DataTable{}
	for _, cmd := range commands {
		result.AddRow(cmd)
	}
	return result
}
