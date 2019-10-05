package steps

import (
	"fmt"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/Originate/git-town/test"
	"github.com/Originate/git-town/test/gherkintools"
)

// RunSteps defines Gherkin step implementations around running things in subshells.
func RunSteps(suite *godog.Suite, fs *FeatureState) {
	suite.Step(`^I run "([^"]*)"$`, func(command string) error {
		fs.activeScenarioState.lastRunOutput, fs.activeScenarioState.lastRunErr = fs.activeScenarioState.gitEnvironment.DeveloperRepo.RunString(command)
		return nil
	})

	suite.Step(`^it runs the commands$`, func(table *gherkin.DataTable) error {
		commands := test.GitCommandsInGitTownOutput(fs.activeScenarioState.lastRunOutput)
		return gherkintools.EnsureStringSliceMatchesTable(commands, table)
	})

	suite.Step(`^it runs no commands$`, func() error {
		commands := test.GitCommandsInGitTownOutput(fs.activeScenarioState.lastRunOutput)
		if len(commands) > 0 {
			for _, command := range commands {
				fmt.Println(command)
			}
			return fmt.Errorf("expected no commands but found %d commands", len(commands))
		}
		return nil
	})
}
