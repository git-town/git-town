package main_test

import (
	"os"
	"runtime"
	"testing"

	"github.com/cucumber/godog"
	"github.com/git-town/git-town/v14/test/cucumber"
	"github.com/spf13/pflag"
)

// func FeatureContext(suite *godog.Suite) {
// 	// The current Godog implementation only provides a FeatureContext,
// 	// no SuiteContext nor ScenarioContext.
// 	// Hence we have to register the scenario state here (and reuse it for all scenarios in a feature)
// 	// and register the steps here.
// 	// It is initialized in SuiteSteps.BeforeScenario.
// 	state := cucumber.ScenarioState{}
// 	cucumber.Steps(suite, &state)
// }

//nolint:paralleltest
func TestGodog(t *testing.T) {
	var options = godog.Options{
		// DefaultContext: ,
		Format:        "progress",
		StopOnFailure: true,
		// Strict:        true,
	}
	godog.BindCommandLineFlags("godog.", &options)
	pflag.Parse()
	options.Paths = pflag.Args()
	if len(options.Paths) > 0 {
		options.Format = "pretty"
	}
	if runtime.GOOS == "windows" {
		options.Tags = "~@skipWindows"
		options.Concurrency = runtime.NumCPU()
	} else {
		options.Concurrency = runtime.NumCPU() * 4
	}
	if os.Getenv("smoke") != "" {
		options.Tags = "@smoke"
	}
	if os.Getenv("cukethis") != "" {
		options.Tags = "@this"
	}
	suite := godog.TestSuite{
		Name:                 "godogs",
		TestSuiteInitializer: cucumber.InitializeSuite,
		ScenarioInitializer:  cucumber.InitializeScenario,
		Options:              &options,
	}
	status := suite.Run()
	if status > 0 {
		t.FailNow()
	}
}
