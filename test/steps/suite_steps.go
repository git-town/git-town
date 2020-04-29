package steps

import (
	"fmt"
	"io/ioutil"
	"log"
	"sync"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/gherkin"
	"github.com/git-town/git-town/test"
)

// beforeSuiteMux ensures that we run BeforeSuite only once globally.
var beforeSuiteMux sync.Mutex

// the global GitManager instance
var gitManager *test.GitManager

// SuiteSteps defines global lifecycle step implementations for Cucumber.
func SuiteSteps(suite *godog.Suite, fs *FeatureState) {
	suite.BeforeFeature(func(feature *gherkin.Feature) {
		if hasFeatureTag(feature, "@debug") {
			test.Debug = true
		}
	})

	suite.BeforeScenario(func(args interface{}) {
		// create a GitEnvironment for the scenario
		gitEnvironment, err := gitManager.CreateScenarioEnvironment(scenarioName(args))
		if err != nil {
			log.Fatalf("cannot create environment for scenario %q: %s", scenarioName(args), err)
		}
		fs.activeScenarioState = scenarioState{gitEnvironment: gitEnvironment}
		if hasScenarioTag(args, "@debug") {
			test.Debug = true
		}
	})

	suite.BeforeSuite(func() {
		// NOTE: we want to create only one global GitManager instance with one global memoized environment.
		beforeSuiteMux.Lock()
		defer beforeSuiteMux.Unlock()
		if gitManager == nil {
			baseDir, err := ioutil.TempDir("", "")
			if err != nil {
				log.Fatalf("cannot create base directory for feature specs: %s", err)
			}
			gitManager = test.NewGitManager(baseDir)
			err = gitManager.CreateMemoizedEnvironment()
			if err != nil {
				log.Fatalf("Cannot create memoized environment: %s", err)
			}
		}
	})

	suite.AfterScenario(func(args interface{}, e error) {
		if e == nil {
			err := fs.activeScenarioState.gitEnvironment.Remove()
			if err != nil {
				log.Fatalf("error removing the Git environment after scenario %q: %v", scenarioName(args), err)
			}
		} else {
			fmt.Printf("failed scenario, investigate state in %q\n", fs.activeScenarioState.gitEnvironment.Dir)
		}
	})
}

// hasFeatureTag indicates whether the given feature has a tag with the given name.
func hasFeatureTag(feature *gherkin.Feature, name string) bool {
	for _, tag := range feature.Tags {
		if tag.Name == name {
			return true
		}
	}
	return false
}

// hasScenarioTag indicates whether the given scenario or scenario outline has a tag of the given name.
func hasScenarioTag(args interface{}, name string) bool {
	for _, tag := range scenarioTags(args) {
		if tag.Name == name {
			return true
		}
	}
	return false
}

// scenarioName returns the name of the given Scenario or ScenarioOutline
func scenarioName(args interface{}) string {
	scenario, ok := args.(*gherkin.Scenario)
	if ok {
		return scenario.Name
	}
	scenarioOutline, ok := args.(*gherkin.ScenarioOutline)
	if ok {
		return scenarioOutline.Name
	}
	panic("unknown scenario type")
}

func scenarioTags(args interface{}) []*gherkin.Tag {
	switch typed := args.(type) {
	case *gherkin.Scenario:
		return typed.Tags
	case *gherkin.ScenarioOutline:
		return typed.Tags
	default:
		panic("unknown scenario type")
	}
}
