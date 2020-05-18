package main_test

import (
	"runtime"
	"testing"

	"github.com/cucumber/godog"
	"github.com/git-town/git-town/test/steps"
)

// nolint:deadcode,unused
func FeatureContext(suite *godog.Suite) {
	// The current Godog implementation only provides a FeatureContext,
	// no SuiteContext nor ScenarioContext.
	// Hence we have to register the scenario state here (and reuse it for all scenarios in a feature)
	// and register the steps here.
	// It is initialized in SuiteSteps.BeforeScenario.
	state := &steps.ScenarioState{}
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

func TestGodog(t *testing.T) {
	status := godog.RunWithOptions("godog", func(s *godog.Suite) {
		FeatureContext(s)
	}, godog.Options{
		Format:      "progress",
		Concurrency: runtime.NumCPU(),
		Strict:      true,
		Paths:       []string{"features/"},
	})
	if status > 0 {
		t.FailNow()
	}
}
