package steps

import (
	"fmt"

	"github.com/cucumber/godog"
	"github.com/cucumber/messages-go/v10"
	"github.com/git-town/git-town/test"
	"github.com/git-town/git-town/test/helpers"
)

// CommitSteps defines Cucumber step implementations around configuration.
func CommitSteps(suite *godog.Suite, fs *FeatureState) {
	suite.Step(`^my repository is left with my original commits$`, func() error {
		return compareExistingCommits(fs, fs.activeScenarioState.originalCommitTable)
	})

	suite.Step(`^my repository now has the following commits$`, func(table *messages.PickleStepArgument_PickleTable) error {
		return compareExistingCommits(fs, table)
	})

	suite.Step(`^the following commits exist in my repository$`, func(table *messages.PickleStepArgument_PickleTable) error {
		fs.activeScenarioState.originalCommitTable = table
		commits, err := test.FromGherkinTable(table)
		if err != nil {
			return fmt.Errorf("cannot parse Gherkin table: %w", err)
		}
		return fs.activeScenarioState.gitEnvironment.CreateCommits(commits)
	})
}

// compareExistingCommits compares the commits in the Git environment of the given FeatureState
// against the given Gherkin table.
func compareExistingCommits(fs *FeatureState, table *messages.PickleStepArgument_PickleTable) error {
	fields := helpers.TableFields(table)
	commitTable, err := fs.activeScenarioState.gitEnvironment.CommitTable(fields)
	if err != nil {
		return fmt.Errorf("cannot determine commits in the developer repo: %w", err)
	}
	diff, errorCount := commitTable.EqualGherkin(table)
	if errorCount != 0 {
		fmt.Printf("\nERROR! Found %d differences in the existing commits\n\n", errorCount)
		fmt.Println(diff)
		return fmt.Errorf("mismatching commits found, see diff above")
	}
	return nil
}
