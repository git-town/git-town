package main_test

import (
	"github.com/DATA-DOG/godog"
	"github.com/Originate/git-town/test/steps"
)

// nolint:deadcode
func FeatureContext(s *godog.Suite) {
	steps.SuiteSteps(s)
	steps.ConfigurationSteps(s)
	steps.InstallationSteps(s)
	steps.PrintSteps(s)
	steps.RunSteps(s)
	steps.WorkspaceSteps(s)
}
