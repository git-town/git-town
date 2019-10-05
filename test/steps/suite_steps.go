package steps

import (
	"io/ioutil"
	"log"
	"sync"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/Originate/git-town/test"
	"github.com/Originate/git-town/test/helpers"
	"github.com/iancoleman/strcase"
)

// mux ensures that we run BeforeSuite only once globally.
var mux sync.Mutex

// SuiteSteps defines global lifecycle step implementations for Cucumber.
func SuiteSteps(s *godog.Suite, state *FeatureState) {
	s.BeforeSuite(func() {
		// NOTE: we want to create only one global GitManager instance with one global memoized environment.
		mux.Lock()
		defer mux.Unlock()
		if gitManager == nil {
			baseDir, err := ioutil.TempDir("", "")
			if err != nil {
				log.Fatalf("cannot create base directory: %s", err)
			}
			gitManager = test.NewGitManager(baseDir)
			err = gitManager.CreateMemoizedEnvironment()
			if err != nil {
				log.Fatalf("Cannot create memoized environment: %s", err)
			}
		}
	})

	s.BeforeScenario(state.beforeScenario)
	s.AfterScenario(state.afterScenario)
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
	panic("unknown type")
}

func (state *FeatureState) beforeScenario(args interface{}) {
	// create a GitEnvironment for the scenario
	environmentName := strcase.ToSnake(scenarioName(args)) + "_" + helpers.RandomNumberString(10)
	var err error
	state.gitEnvironment, err = gitManager.CreateScenarioEnvironment(environmentName)
	if err != nil {
		log.Fatalf("cannot create environment for scenario '%s': %s", environmentName, err)
	}
}

func (state *FeatureState) afterScenario(args interface{}, e error) {
	// remove the GitEnvironment of the scenario
	err := state.gitEnvironment.Remove()
	if err != nil {
		log.Fatalf("error in afterScenario: %v", err)
	}
}
