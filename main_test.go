package main_test

import (
	"os"
	"testing"

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

func TestMain(m *testing.M) {
	status := godog.RunWithOptions("godog", func(s *godog.Suite) {
		FeatureContext(s)
	}, godog.Options{
		Format: "progress",
		Paths:  []string{"features/git-town-append/offline.feature:18"},
	})
	if st := m.Run(); st > status {
		status = st
	}
	os.Exit(status)
}
