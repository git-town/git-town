package steps

import (
	"github.com/cucumber/godog"
)

// RepoSteps defines Gherkin step implementations around running things in subshells.
func RepoSteps(suite *godog.Suite, state *ScenarioState) {
	suite.Step(`^my repo has an upstream repo$`, func() error {
		return state.gitEnv.AddUpstream()
	})

	suite.Step(`^my repository knows about the remote branch$`, func() error {
		return state.gitEnv.DevRepo.Fetch()
	})
}
