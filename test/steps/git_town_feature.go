package steps

import "github.com/Originate/git-town/test"

// gitTownFeature contains state that is shared by step implementations.
type GitTownFeature struct {

	// the GitEnvironment used in the current scenario
	gitEnvironment *test.GitEnvironment

	// the result of the last run of Git Town
	lastRunOutput string

	// the error of the last run of Git Town
	lastRunErr error
}
