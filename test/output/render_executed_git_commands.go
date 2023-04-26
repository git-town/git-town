package output

import (
	"github.com/cucumber/messages-go/v10"
	"github.com/git-town/git-town/v8/test/gherkin"
)

// RenderExecutedGitCommands provides the textual Gherkin table representation of the given executed Git commands.
// The DataTable table matches the structure of the given Gherkin table.
func RenderExecutedGitCommands(commands []ExecutedGitCommand, table *messages.PickleStepArgument_PickleTable) gherkin.DataTable {
	tableHasBranches := table.Rows[0].Cells[0].Value == "BRANCH"
	result := gherkin.DataTable{}
	if tableHasBranches {
		result.AddRow("BRANCH", "COMMAND")
	} else {
		result.AddRow("COMMAND")
	}
	lastBranch := ""
	for _, cmd := range commands {
		if tableHasBranches {
			switch {
			case cmd.Branch == lastBranch:
				result.AddRow("", cmd.Command)
			case cmd.Branch == "":
				result.AddRow("<none>", cmd.Command)
			default:
				result.AddRow(cmd.Branch, cmd.Command)
			}
		} else {
			result.AddRow(cmd.Command)
		}
		lastBranch = cmd.Branch
	}
	return result
}
