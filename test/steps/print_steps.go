package steps

import (
	"fmt"
	"strings"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
)

// PrintSteps defines Gherkin steps around printing things to the terminal.
func PrintSteps(suite *godog.Suite, fs *FeatureState) {
	suite.Step("^it prints$", fs.itPrints)
	suite.Step("^it does not print \"([^\"]*)\"$", fs.itDoesNotPrint)
	suite.Step(`^it prints the error:$`, fs.itPrintsTheError)
}

func (fs *FeatureState) itPrints(expected *gherkin.DocString) error {
	if !strings.Contains(fs.activeScenarioState.lastRunOutput, expected.Content) {
		return fmt.Errorf(`text not found: %s`, expected.Content)
	}
	return nil
}

func (fs *FeatureState) itDoesNotPrint(text string) error {
	if strings.Contains(fs.activeScenarioState.lastRunOutput, text) {
		return fmt.Errorf(`text found: %s`, text)
	}
	return nil
}

func (fs *FeatureState) itPrintsTheError(expected *gherkin.DocString) error {
	if !strings.Contains(fs.activeScenarioState.lastRunOutput, expected.Content) {
		return fmt.Errorf("text not found: %s\n\nactual text:\n%s", expected.Content, fs.activeScenarioState.lastRunOutput)
	}
	if fs.activeScenarioState.lastRunErr == nil {
		return fmt.Errorf("expected error")
	}
	return nil
}
