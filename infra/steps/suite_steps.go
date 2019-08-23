package steps

import (
	"io/ioutil"
	"log"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/Originate/git-town/infra"
	"github.com/dchest/uniuri"
	"github.com/iancoleman/strcase"
)

func SuiteSteps(s *godog.Suite) {
	s.BeforeSuite(func() {

		// create the directory to put the GitEnvironments ino
		baseDir, err := ioutil.TempDir("", "")
		if err != nil {
			log.Fatalf("cannot create base directory: %s", err)
		}

		// create the GitManager
		gitManager = infra.NewGitManager(baseDir)

		// create the memoized environment
		err = gitManager.CreateMemoizedEnvironment()
		if err != nil {
			log.Fatalf("Cannot create memoized environment: %s", err)
		}
	})

	s.BeforeScenario(func(args interface{}) {
		// create a GitEnvironment for the scenario
		environmentName := strcase.ToSnake(scenarioName(args)) + "_" + string(uniuri.NewLen(10))
		var err error
		gitEnvironment, err = gitManager.CreateScenarioEnvironment(environmentName)
		if err != nil {
			log.Fatalf("cannot create environment for scenario '%s': %s", environmentName, err)
		}
	})

	s.AfterScenario(func(args interface{}, err error) {
		// TODO: delete scenario environment
	})
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
