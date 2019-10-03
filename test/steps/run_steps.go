package steps

import (
	"fmt"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/Originate/git-town/test"
	"github.com/Originate/git-town/test/cucumber"
)

// RunSteps defines Gherkin step implementations around running things in subshells.
func RunSteps(s *godog.Suite) {
	s.Step(`^I run "([^"]*)"$`, func(command string) error {
		lastRunOutput, lastRunErr = gitEnvironment.DeveloperRepo.RunString(command)
		return nil
	})

	s.Step(`^it runs the commands$`,
		func(table *gherkin.DataTable) error {
			commands := test.GitCommandsInGitTownOutput(lastRunOutput)
			return cucumber.AssertStringSliceMatchesTable(commands, table)
		})

	s.Step(`^it runs no commands$`,
		func() error {
			commands := test.GitCommandsInGitTownOutput(lastRunOutput)
			if len(commands) > 0 {
				for _, command := range commands {
					fmt.Println(command)
				}
				return fmt.Errorf("expected no commands but found %d commands", len(commands))
			}
			return nil
		})

}
