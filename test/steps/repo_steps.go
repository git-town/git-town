package steps

import (
	"github.com/cucumber/godog"
)

// RepoSteps defines Gherkin step implementations around running things in subshells.
func RepoSteps(suite *godog.Suite, fs *FeatureState) {
	suite.Step(`^my repo has an upstream repo$`, func() error {
		return fs.state.gitEnv.AddUpstream()
	})

	suite.Step(`^my repository knows about the remote branch$`, func() error {
		return fs.state.gitEnv.DeveloperRepo.Fetch()
	})
}
