package test

import (
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/Originate/git-town/test/gherkintools"
)

// RenderExecutedGitCommands provides the textual Gherkin table representation of the given executed Git commands.
func RenderExecutedGitCommands(commands []ExecutedGitCommand) string {
	renderer := gherkintools.TableRenderer{}
	renderer.AddLine("BRANCH", "COMMAND")
	lastBranch := ""
	for _, cmd := range commands {
		if cmd.Branch == lastBranch {
			renderer.AddLine("", cmd.Command)
		} else {
			renderer.AddLine(cmd.Branch, cmd.Command)
		}
		lastBranch = cmd.Branch
	}
	return renderer.String()
}

// RenderTable provides the textual Gherkin representation of the given Gherkin table.
func RenderTable(table *gherkin.DataTable) string {
	renderer := gherkintools.TableRenderer{}
	for _, row := range table.Rows {
		values := []string{}
		for _, cell := range row.Cells {
			values = append(values, cell.Value)
		}
		renderer.AddLine(values...)
	}
	return renderer.String()
}
