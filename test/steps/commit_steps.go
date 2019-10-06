package steps

import (
	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
)

// CommitSteps defines Cucumber step implementations around configuration.
func CommitSteps(suite *godog.Suite, fs *FeatureState) {
	suite.Step(`^the following commits exist in my repository$`, func(table *gherkin.DataTable) error {
		return fs.activeScenarioState.gitEnvironment.CreateCommits(table)
	})
}
