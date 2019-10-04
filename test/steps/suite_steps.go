package steps

import (
	"io/ioutil"
	"log"
	"sync"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/Originate/git-town/test"
	"github.com/dchest/uniuri"
	"github.com/iancoleman/strcase"
)

// mux ensures that we run BeforeSuite only once globally.
var mux sync.Mutex

// SuiteSteps provides global lifecycle step implementations for Cucumber.
func SuiteSteps(s *godog.Suite, gtf *GitTownFeature) {
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

	s.BeforeScenario(gtf.beforeScenario)
	s.AfterScenario(gtf.afterScenario)
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

func (gtf *GitTownFeature) beforeScenario(args interface{}) {
	// create a GitEnvironment for the scenario
	environmentName := strcase.ToSnake(scenarioName(args)) + "_" + string(uniuri.NewLen(10))
	var err error
	gtf.gitEnvironment, err = gitManager.CreateScenarioEnvironment(environmentName)
	if err != nil {
		log.Fatalf("cannot create environment for scenario '%s': %s", environmentName, err)
	}
}

func (gtf *GitTownFeature) afterScenario(args interface{}, err error) {
	// TODO: delete scenario environment
}
