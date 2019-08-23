package main

/*
Test setup:
- by default, each scenario runs in a directory called "developer"
	that has a "main" branch and a valid Git Town configuration
- at script startup, it creates a memoized repo with that setup
- before each scenario, it copies that memoized repo over into the "developer" repo
*/

import (
	"github.com/DATA-DOG/godog"
	"github.com/Originate/git-town/infra/steps"
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
