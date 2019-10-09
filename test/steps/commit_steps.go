package steps

import (
	"fmt"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/Originate/git-town/test/gherkintools"
	"github.com/pkg/errors"
)

// CommitSteps defines Cucumber step implementations around configuration.
func CommitSteps(suite *godog.Suite, fs *FeatureState) {
	suite.Step(`^the following commits exist in my repository$`, func(table *gherkin.DataTable) error {
		return fs.activeScenarioState.gitEnvironment.CreateCommits(table)
	})

	suite.Step(`^my repository now has the following commits$`, func(table *gherkin.DataTable) error {
		fields := gherkintools.TableFields(table)
		commits, err := fs.activeScenarioState.gitEnvironment.Commits(fields)
		if err != nil {
			return errors.Wrap(err, "cannot determine commits in the developer repo")
		}
		diff, errorCount := commits.Equal(table)
		if errorCount != 0 {
			fmt.Println(diff)
			return fmt.Errorf("found %d differences", errorCount)
		}
		return nil
	})
}
