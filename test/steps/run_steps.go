package steps

import (
	"fmt"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/Originate/git-town/test"
	"github.com/Originate/git-town/test/cucumber"
)

// RunSteps defines Gherkin step implementations around running things in subshells.
func RunSteps(s *godog.Suite, gtf *GitTownFeature) {
	s.Step(`^I run "([^"]*)"$`, gtf.iRun)
	s.Step(`^it runs the commands$`, gtf.itRunsTheCommands)
	s.Step(`^it runs no commands$`, gtf.itRunsNoCommands)
}

func (gtf *GitTownFeature) iRun(command string) error {
	gtf.lastRunOutput, gtf.lastRunErr = gtf.gitEnvironment.DeveloperRepo.RunString(command)
	return nil
}

func (gtf *GitTownFeature) itRunsTheCommands(table *gherkin.DataTable) error {
	commands := test.GitCommandsInGitTownOutput(gtf.lastRunOutput)
	return cucumber.AssertStringSliceMatchesTable(commands, table)
}

func (gtf *GitTownFeature) itRunsNoCommands() error {
	commands := test.GitCommandsInGitTownOutput(gtf.lastRunOutput)
	if len(commands) > 0 {
		for _, command := range commands {
			fmt.Println(command)
		}
		return fmt.Errorf("expected no commands but found %d commands", len(commands))
	}
	return nil
}
