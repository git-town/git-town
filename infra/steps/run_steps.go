package steps

import (
	"fmt"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/Originate/git-town/infra"
)

func RunSteps(s *godog.Suite) {
	s.Step(`^I run "([^"]*)"$`, func(command string) error {
		lastRunResult = gitEnvironment.DeveloperRepo.RunString(command)
		return nil
	})

	s.Step(`^it runs the commands$`,
		func(table *gherkin.DataTable) error {
			commands := infra.CommandsInOutput(lastRunResult.Output)
			return infra.AssertStringSliceMatchesTable(commands, table)
		})

	s.Step(`^it runs no commands$`,
		func() error {
			commands := infra.CommandsInOutput(lastRunResult.Output)
			if len(commands) > 0 {
				for _, command := range commands {
					fmt.Println(command)
				}
				return fmt.Errorf("expected no commands but found %d commands", len(commands))
			}
			return nil
		})

}
