package main_test

import (
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/cucumber/godog"
	"github.com/git-town/git-town/v19/internal/test/cucumber"
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
	flagThis := len(os.Getenv("cukethis")) > 0
	flagSmoke := len(os.Getenv("smoke")) > 0
	flagMessyOutput := os.Getenv("messyoutput")
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
	// options.Format = "pretty"
	if runtime.GOOS == "windows" {
		options.Tags = "~@skipWindows"
		options.Concurrency = runtime.NumCPU()
	} else {
		options.Concurrency = runtime.NumCPU() * 4
	}
	switch flagMessyOutput {
	case "":
	case "0":
		options.Tags = "~@messyoutput"
	case "1":
		options.Tags = "@messyoutput"
	default:
		panic("unknown value for messyoutput")
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
