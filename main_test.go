package main_test

import (
	"os"
	"testing"

	"github.com/cucumber/godog"
	"github.com/git-town/git-town/test/steps"
)

// nolint:deadcode,unused
func FeatureContext(suite *godog.Suite) {
	state := &steps.FeatureState{}
	steps.SuiteSteps(suite, state)
	steps.AutocompletionSteps(suite, state)
	steps.BranchSteps(suite, state)
	steps.CommitSteps(suite, state)
	steps.ConfigurationSteps(suite, state)
	steps.ConflictSteps(suite, state)
	steps.DebugSteps(suite, state)
	steps.FileSteps(suite, state)
	steps.FolderSteps(suite, state)
	steps.GitTownSteps(suite, state)
	steps.InstallationSteps(suite, state)
	steps.MergeSteps(suite, state)
	steps.OfflineSteps(suite, state)
	steps.OriginSteps(suite, state)
	steps.PrintSteps(suite, state)
	steps.RebaseSteps(suite, state)
	steps.RepoSteps(suite, state)
	steps.RunSteps(suite, state)
	steps.WorkspaceSteps(suite, state)
	steps.MergeSteps(suite, state)
	steps.TagSteps(suite, state)
	steps.CoworkerSteps(suite, state)
}

func TestMain(m *testing.M) {
	status := godog.RunWithOptions("godog", func(s *godog.Suite) {
		FeatureContext(s)
	}, godog.Options{
		Format: "pretty",
		Paths:  []string{"features/git-town-append/on-perennial-branch.feature"},
	})
	if st := m.Run(); st > status {
		status = st
	}
	os.Exit(status)
}
