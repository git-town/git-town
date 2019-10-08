package test

import (
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/Originate/git-town/test/gherkintools"
)

// RenderExecutedGitCommands provides the textual Gherkin table representation of the given executed Git commands.
func RenderExecutedGitCommands(commands []ExecutedGitCommand) gherkintools.Mortadella {
	morta := gherkintools.Mortadella{}
	morta.AddRow("BRANCH", "COMMAND")
	lastBranch := ""
	for _, cmd := range commands {
		if cmd.Branch == lastBranch {
			morta.AddRow("", cmd.Command)
		} else {
			morta.AddRow(cmd.Branch, cmd.Command)
		}
		lastBranch = cmd.Branch
	}
	return morta
}

// RenderTable provides the textual Gherkin representation of the given Gherkin table.
func RenderTable(table *gherkin.DataTable) string {
	morta := gherkintools.Mortadella{}
	for _, row := range table.Rows {
		values := []string{}
		for _, cell := range row.Cells {
			values = append(values, cell.Value)
		}
		morta.AddRow(values...)
	}
	return morta.String()
}
