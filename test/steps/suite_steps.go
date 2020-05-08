package steps

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"sync"

	"github.com/cucumber/godog"
	"github.com/cucumber/messages-go/v10"
	"github.com/git-town/git-town/test"
)

// beforeSuiteMux ensures that we run BeforeSuite only once globally.
var beforeSuiteMux sync.Mutex

// the global GitManager instance
var gitManager *test.GitManager

// SuiteSteps defines global lifecycle step implementations for Cucumber.
func SuiteSteps(suite *godog.Suite, fs *FeatureState) {
	suite.BeforeScenario(func(scenario *messages.Pickle) {
		// create a GitEnvironment for the scenario
		gitEnvironment, err := gitManager.CreateScenarioEnvironment(scenarioName(scenario))
		if err != nil {
			log.Fatalf("cannot create environment for scenario %q: %s", scenarioName(scenario), err)
		}
		fs.activeScenarioState = scenarioState{gitEnvironment: gitEnvironment}
		if hasTag(scenario, "@debug") {
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
			// Evaluate symlinks as Mac temp dir is symlinked
			evalBaseDir, err := filepath.EvalSymlinks(baseDir)
			if err != nil {
				log.Fatalf("cannot evaluate symlinks of base directory for feature specs: %s", err)
			}
			gitManager = test.NewGitManager(evalBaseDir)
			err = gitManager.CreateMemoizedEnvironment()
			if err != nil {
				log.Fatalf("Cannot create memoized environment: %s", err)
			}
		}
	})

	suite.AfterScenario(func(scenario *messages.Pickle, e error) {
		if e == nil {
			err := fs.activeScenarioState.gitEnvironment.Remove()
			if err != nil {
				log.Fatalf("error removing the Git environment after scenario %q: %v", scenarioName(scenario), err)
			}
		} else {
			fmt.Printf("failed scenario, investigate state in %q\n", fs.activeScenarioState.gitEnvironment.Dir)
		}
	})
}

// hasTag indicates whether the given feature has a tag with the given name.
func hasTag(scenario *messages.Pickle, name string) bool {
	for _, tag := range scenario.GetTags() {
		if tag.Name == name {
			return true
		}
	}
	return false
}

// scenarioName returns the name of the given Scenario or ScenarioOutline
func scenarioName(args *messages.Pickle) string {
	return args.GetName()
}
