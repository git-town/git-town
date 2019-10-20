package steps

import (
	"fmt"
	"io/ioutil"
	"log"
	"sync"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/Originate/git-town/test"
)

// beforeSuiteMux ensures that we run BeforeSuite only once globally.
var beforeSuiteMux sync.Mutex

// the global GitManager instance
var gitManager *test.GitManager

// SuiteSteps defines global lifecycle step implementations for Cucumber.
func SuiteSteps(suite *godog.Suite, fs *FeatureState) {
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

	suite.BeforeScenario(func(args interface{}) {
		gitEnvironment, err := gitManager.CreateScenarioEnvironment(scenarioName(args))
		if err != nil {
			log.Fatalf("cannot create environment for scenario %q: %s", scenarioName(args), err)
		}
		fs.activeScenarioState = scenarioState{gitEnvironment: gitEnvironment}
	})

	suite.AfterScenario(func(args interface{}, e error) {
		if e != nil {
			fmt.Printf("failed scenario, investigate state in %q\n", fs.activeScenarioState.gitEnvironment.Dir)
		} else {
			err := fs.activeScenarioState.gitEnvironment.Remove()
			if err != nil {
				log.Fatalf("error removing the Git environment after scenario %q: %v", scenarioName(args), err)
			}
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
	scenario, ok := args.(*gherkin.Scenario)
	if ok {
		return scenario.Tags
	}
	scenarioOutline, ok := args.(*gherkin.ScenarioOutline)
	if ok {
		return scenarioOutline.Tags
	}
	panic("unknown scenario type")
}
