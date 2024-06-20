package main_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"runtime"
	"testing"

	"github.com/cucumber/godog"
	"github.com/git-town/git-town/v14/test/cucumber"
	"github.com/spf13/pflag"
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
	var options = godog.Options{
		// DefaultContext: ,
		Format:        "progress",
		StopOnFailure: true,
		Strict:        true,
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
		Name: "godogs",
		// TestSuiteInitializer: InitializeTestSuite,
		ScenarioInitializer: InitializeScenario,
		Options:             &options,
	}
	status := suite.Run()
	os.Exit(status)
}

func InitializeSuite(ctx *godog.TestSuiteContext) {
	ctx.BeforeSuite(func() {
		fmt.Println("BEFORE SUITE")
	})
	ctx.AfterSuite(func() {
		fmt.Println("AFTER SUITE")
	})
	sc := ctx.ScenarioContext()

	sc.Given(`^there are (\d+) godogs$`, func(ctx context.Context, available int) (context.Context, error) {
		return context.WithValue(ctx, keyGodogs, available), nil
	})

	sc.When(`^I eat (\d+)$`, func(ctx context.Context, num int) (context.Context, error) {
		available, ok := ctx.Value(keyGodogs).(int)
		if !ok {
			return ctx, errors.New("there are no godogs available")
		}
		if available < num {
			return ctx, fmt.Errorf("you cannot eat %d godogs, there are %d available", num, available)
		}
		available -= num
		return context.WithValue(ctx, keyGodogs, available), nil
	})

	sc.Then(`^there should be (\d+) remaining$`, func(ctx context.Context, remaining int) error {
		available, has := ctx.Value(keyGodogs).(int)
		if !has {
			return errors.New("there are no godogs available")
		}
		if available != remaining {
			return fmt.Errorf("expected %d godogs to be remaining, but there is %d", remaining, available)
		}
		return nil
	})
}
