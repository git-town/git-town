package steps

import (
	"fmt"
	"strings"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/Originate/git-town/src/command"
	"github.com/Originate/git-town/test"
)

// RunSteps defines Gherkin step implementations around running things in subshells.
func RunSteps(suite *godog.Suite, fs *FeatureState) {
	suite.Step(`^I run "([^"]+)"$`, func(command string) error {
		fs.activeScenarioState.lastRunResult, fs.activeScenarioState.lastRunErr = fs.activeScenarioState.gitEnvironment.DeveloperRepo.RunString(command)
		return nil
	})

	suite.Step(`^I run "([^"]+)" and answer the prompts:$`, func(cmd string, input *gherkin.DataTable) error {
		fs.activeScenarioState.lastRunResult, fs.activeScenarioState.lastRunErr = fs.activeScenarioState.gitEnvironment.DeveloperRepo.RunStringWith(cmd, command.Options{Input: tableToInput(input)})
		return nil
	})

	suite.Step(`^I run "([^"]+)" in the "([^"]+)" folder$`, func(cmd, folderName string) error {
		fs.activeScenarioState.lastRunResult, fs.activeScenarioState.lastRunErr = fs.activeScenarioState.gitEnvironment.DeveloperRepo.RunStringWith(cmd, command.Options{Dir: folderName})
		return nil
	})

	suite.Step(`^it runs no commands$`, func() error {
		commands := test.GitCommandsInGitTownOutput(fs.activeScenarioState.lastRunResult.Output())
		if len(commands) > 0 {
			for _, command := range commands {
				fmt.Println(command)
			}
			return fmt.Errorf("expected no commands but found %d commands", len(commands))
		}
		return nil
	})

	suite.Step(`^it runs the commands$`, func(input *gherkin.DataTable) error {
		commands := test.GitCommandsInGitTownOutput(fs.activeScenarioState.lastRunResult.Output())
		table := test.RenderExecutedGitCommands(commands, input)
		dataTable := test.FromGherkin(input)
		expanded := dataTable.Expand(fs.activeScenarioState.gitEnvironment.DeveloperRepo.Dir)
		diff, errorCount := table.EqualDataTable(expanded)
		if errorCount != 0 {
			fmt.Printf("\nERROR! Found %d differences in the commands run\n\n", errorCount)
			fmt.Println(diff)
			return fmt.Errorf("mismatching commands run, see diff above")
		}
		return nil
	})
}

func tableToInput(table *gherkin.DataTable) []string {
	var result []string
	for i := 1; i < len(table.Rows); i++ {
		row := table.Rows[i]
		answer := row.Cells[1].Value
		answer = strings.ReplaceAll(answer, "[ENTER]", "\n")
		answer = strings.ReplaceAll(answer, "[DOWN]", "\x1b[B")
		answer = strings.ReplaceAll(answer, "[UP]", "\x1b[A")
		answer = strings.ReplaceAll(answer, "[SPACE]", " ")
		result = append(result, answer)
	}
	return result
}
