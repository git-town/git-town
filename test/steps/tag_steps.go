package steps

import (
	"fmt"

	"github.com/cucumber/godog"
	"github.com/cucumber/messages-go/v10"
)

// TagSteps defines Gherkin step implementations around merges.
func TagSteps(suite *godog.Suite, state *ScenarioState) {
	suite.Step(`^my repo has the following tags$`, func(table *messages.PickleStepArgument_PickleTable) error {
		return state.gitEnv.CreateTags(table)
	})

	suite.Step(`^my repo now has the following tags$`, func(table *messages.PickleStepArgument_PickleTable) error {
		tagTable, err := state.gitEnv.TagTable()
		if err != nil {
			return err
		}
		diff, errorCount := tagTable.EqualGherkin(table)
		if errorCount != 0 {
			fmt.Printf("\nERROR! Found %d differences in the existing tags\n\n", errorCount)
			fmt.Println(diff)
			return fmt.Errorf("mismatching tags found, see diff above")
		}
		return nil
	})

	suite.Step(`^my repo has a remote tag "([^"]+)" that is not on a branch$`, func(name string) error {
		return state.gitEnv.OriginRepo.CreateStandaloneTag(name)
	})
}
