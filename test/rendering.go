package test

import (
	"github.com/cucumber/godog/gherkin"
)

// RenderExecutedGitCommands provides the textual Gherkin table representation of the given executed Git commands.
// The DataTable table matches the structure of the given Gherkin table.
func RenderExecutedGitCommands(commands []ExecutedGitCommand, table *gherkin.DataTable) DataTable {
	tableHasBranches := table.Rows[0].Cells[0].Value == "BRANCH"
	result := DataTable{}
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

// RenderTable provides the textual Gherkin representation of the given Gherkin table.
func RenderTable(table *gherkin.DataTable) string {
	result := DataTable{}
	for _, row := range table.Rows {
		values := []string{}
		for _, cell := range row.Cells {
			values = append(values, cell.Value)
		}
		result.AddRow(values...)
	}
	return result.String()
}
