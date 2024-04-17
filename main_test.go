package main_test

import (
	"os"
	"runtime"
	"testing"

	"github.com/cucumber/godog"
	"github.com/git-town/git-town/v14/test/cucumber"
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
	var concurrency int
	if runtime.GOOS == "windows" {
		tags = "~@skipWindows"
		concurrency = 1
	} else {
		concurrency = 4
	}
	if os.Getenv("smoke") != "" {
		tags = "@smoke"
	}
	if os.Getenv("cukethis") != "" {
		tags = "@this"
	}
	status := godog.RunWithOptions("godog", FeatureContext, godog.Options{
		Format:        "progress",
		Concurrency:   runtime.NumCPU() * concurrency,
		StopOnFailure: true,
		Strict:        true,
		Paths:         []string{"features/switch"},
		Tags:          tags,
	})
	if status > 0 {
		t.FailNow()
	}
}
