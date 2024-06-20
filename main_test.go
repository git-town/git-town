package main_test

import (
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/cucumber/godog"
	"github.com/git-town/git-town/v14/test/cucumber"
	"github.com/spf13/pflag"
)

func TestMain(_ *testing.M) {
	options := godog.Options{
		StopOnFailure: true,
		Strict:        true,
	}
	godog.BindCommandLineFlags("godog.", &options)
	pflag.Parse()
	options.Paths = pflag.Args()
	flagThis := os.Getenv("cukethis") != ""
	flagSmoke := os.Getenv("smoke") != ""
	switch {
	case flagThis:
		options.Format = "pretty"
	case len(options.Paths) == 0:
		options.Format = "progress"
	case strings.HasSuffix(options.Paths[0], ".feature"):
		options.Format = "pretty"
	default:
		options.Format = "progress"
	}
	if runtime.GOOS == "windows" {
		options.Tags = "~@skipWindows"
		options.Concurrency = runtime.NumCPU()
	} else {
		options.Concurrency = runtime.NumCPU() * 4
	}
	if flagSmoke {
		options.Tags = "@smoke"
	}
	if flagThis {
		options.Tags = "@this"
	}
	suite := godog.TestSuite{
		Options:              &options,
		ScenarioInitializer:  cucumber.InitializeScenario,
		TestSuiteInitializer: cucumber.InitializeSuite,
	}
	status := suite.Run()
	os.Exit(status)
}
