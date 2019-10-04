package main_test

import (
	"github.com/DATA-DOG/godog"
	"github.com/Originate/git-town/test/steps"
)

// nolint:deadcode,unused
func FeatureContext(s *godog.Suite) {
	gtf := &steps.GitTownFeature{}
	steps.SuiteSteps(s, gtf)
	steps.ConfigurationSteps(s, gtf)
	steps.InstallationSteps(s, gtf)
	steps.PrintSteps(s, gtf)
	steps.RunSteps(s, gtf)
	steps.WorkspaceSteps(s, gtf)
}
