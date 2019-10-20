package main_test

import (
	"github.com/DATA-DOG/godog"
	"github.com/Originate/git-town/test/steps"
)

// nolint:deadcode,unused
func FeatureContext(suite *godog.Suite) {
	state := &steps.FeatureState{}
	steps.SuiteSteps(suite, state)
	steps.BranchSteps(suite, state)
	steps.CommitSteps(suite, state)
	steps.ConfigurationSteps(suite, state)
	steps.GitTownSteps(suite, state)
	steps.InstallationSteps(suite, state)
	steps.PrintSteps(suite, state)
	steps.RunSteps(suite, state)
	steps.WorkspaceSteps(suite, state)
}
