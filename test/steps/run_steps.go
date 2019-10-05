package steps

import (
	"fmt"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/Originate/git-town/test"
	"github.com/Originate/git-town/test/cucumber"
)

// RunSteps defines Gherkin step implementations around running things in subshells.
func RunSteps(s *godog.Suite, state *FeatureState) {
	s.Step(`^I run "([^"]*)"$`, state.iRun)
	s.Step(`^it runs the commands$`, state.itRunsTheCommands)
	s.Step(`^it runs no commands$`, state.itRunsNoCommands)
}

func (state *FeatureState) iRun(command string) error {
	state.lastRunOutput, state.lastRunErr = state.gitEnvironment.DeveloperRepo.RunString(command)
	return nil
}

func (state *FeatureState) itRunsTheCommands(table *gherkin.DataTable) error {
	commands := test.GitCommandsInGitTownOutput(state.lastRunOutput)
	return cucumber.EnsureStringSliceMatchesTable(commands, table)
}

func (state *FeatureState) itRunsNoCommands() error {
	commands := test.GitCommandsInGitTownOutput(state.lastRunOutput)
	if len(commands) > 0 {
		for _, command := range commands {
			fmt.Println(command)
		}
		return fmt.Errorf("expected no commands but found %d commands", len(commands))
	}
	return nil
}
