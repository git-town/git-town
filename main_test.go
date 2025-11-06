package main_test

import (
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/cucumber/godog"
	"github.com/git-town/git-town/v22/internal/test/cucumber"
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
	flagVerbose := len(os.Getenv("verbose")) > 0
	cucumber.CaptureGoldenMode = len(os.Getenv("capturegolden")) > 0
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
	if flagVerbose {
		options.Paths = append(options.Paths, findVerboseFiles()...)
	}
	suite := godog.TestSuite{
		Options:              &options,
		ScenarioInitializer:  cucumber.InitializeScenario,
		TestSuiteInitializer: cucumber.InitializeSuite,
	}
	status := suite.Run()
	os.Exit(status)
}

func findVerboseFiles() []string {
	var result []string
	err := filepath.WalkDir("features", func(path string, dir fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if strings.Contains(dir.Name(), "verbose") {
			result = append(result, path)
		}
		return nil
	})
	if err != nil {
		panic(err.Error())
	}
	if len(result) == 0 {
		panic("no feature files found")
	}
	return result
}
