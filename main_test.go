package main_test

import (
	"runtime"
	"testing"

	"github.com/cucumber/godog"
	"github.com/git-town/git-town/test"
)

// nolint:deadcode,unused
func FeatureContext(suite *godog.Suite) {
	// The current Godog implementation only provides a FeatureContext,
	// no SuiteContext nor ScenarioContext.
	// Hence we have to register the scenario state here (and reuse it for all scenarios in a feature)
	// and register the steps here.
	// It is initialized in SuiteSteps.BeforeScenario.
	state := &test.ScenarioState{}
	test.Steps(suite, state)
}

func TestGodog(t *testing.T) {
	tags := ""
	// noColors := false
	if runtime.GOOS == "windows" {
		tags = "~@skipWindows"
		// noColors = true
	}
	status := godog.RunWithOptions("godog", func(s *godog.Suite) {
		FeatureContext(s)
	}, godog.Options{
		Format:      "progress",
		Concurrency: runtime.NumCPU(),
		Strict:      true,
		// NoColors:    noColors,
		Paths: []string{"features/"},
		Tags:  tags,
	})
	if status > 0 {
		t.FailNow()
	}
}
