package main_test

import (
	"github.com/DATA-DOG/godog"
	"github.com/Originate/git-town/test/steps"
)

// nolint:deadcode,unused
func FeatureContext(s *godog.Suite) {
	state := &steps.FeatureState{}
	steps.SuiteSteps(s, state)
	steps.ConfigurationSteps(s, state)
	steps.InstallationSteps(s, state)
	steps.PrintSteps(s, state)
	steps.RunSteps(s, state)
	steps.WorkspaceSteps(s, state)
}
