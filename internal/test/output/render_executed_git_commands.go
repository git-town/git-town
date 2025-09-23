package output

import (
	"github.com/cucumber/godog"
	"github.com/git-town/git-town/v22/internal/test/datatable"
)

// RenderExecutedGitCommands provides the textual Gherkin table representation of the given executed Git commands.
// The DataTable table matches the structure of the given Gherkin table.
func RenderExecutedGitCommands(commands []ExecutedGitCommand, table *godog.Table) datatable.DataTable {
	tableHasBranches := table.Rows[0].Cells[0].Value == "BRANCH"
	tableHasType := len(table.Rows[0].Cells) > 1 && table.Rows[0].Cells[1].Value == "TYPE"
	result := datatable.DataTable{}
	if tableHasBranches {
		if tableHasType {
			result.AddRow("BRANCH", "TYPE", "COMMAND")
		} else {
			result.AddRow("BRANCH", "COMMAND")
		}
	} else {
		result.AddRow("COMMAND")
	}
	lastBranch := ""
	for _, cmd := range commands {
		if tableHasBranches {
			branch := branchToDisplay(cmd, lastBranch)
			if tableHasType {
				result.AddRow(branch, cmd.CommandType.String(), cmd.Command)
			} else {
				result.AddRow(branch, cmd.Command)
			}
		} else {
			result.AddRow(cmd.Command)
		}
		if cmd.Branch != "" {
			lastBranch = cmd.Branch
		}
	}
	return result
}

func branchToDisplay(cmd ExecutedGitCommand, lastBranch string) string {
	if cmd.Branch == lastBranch {
		return ""
	}
	return cmd.Branch
}
