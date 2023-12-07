package main_test

import (
	"os"
	"runtime"
	"testing"

	"github.com/cucumber/godog"
	"github.com/git-town/git-town/v11/test/cucumber"
)

func FeatureContext(suite *godog.Suite) {
	// The current Godog implementation only provides a FeatureContext,
	// no SuiteContext nor ScenarioContext.
	// Hence we have to register the scenario state here (and reuse it for all scenarios in a feature)
	// and register the steps here.
	// It is initialized in SuiteSteps.BeforeScenario.
	state := cucumber.ScenarioState{}
	cucumber.Steps(suite, &state)
}

//nolint:paralleltest
func TestGodog(t *testing.T) {
	tags := ""
	if os.Getenv("cukethis") != "" {
		tags = "@this"
	}
	if runtime.GOOS == "windows" {
		tags = "~@skipWindows"
	}
	status := godog.RunWithOptions("godog", func(s *godog.Suite) {
		FeatureContext(s)
	}, godog.Options{
		Format:      "progress",
		Concurrency: runtime.NumCPU() * 4,
		Strict:      true,
		Paths:       []string{"features/"},
		Tags:        tags,
	})
	if status > 0 {
		t.FailNow()
	}
}
