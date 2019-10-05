package steps

import (
	"fmt"
	"strings"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
)

// PrintSteps defines Gherkin steps around printing things to the terminal.
func PrintSteps(s *godog.Suite, state *FeatureState) {
	s.Step("^it prints$", state.itPrints)
	s.Step("^it does not print \"([^\"]*)\"$", state.itDoesNotPrint)
	s.Step(`^it prints the error:$`, state.itPrintsTheError)
}

func (state *FeatureState) itPrints(expected *gherkin.DocString) error {
	if !strings.Contains(state.lastRunOutput, expected.Content) {
		return fmt.Errorf(`text not found: %s`, expected.Content)
	}
	return nil
}

func (state *FeatureState) itDoesNotPrint(text string) error {
	if strings.Contains(state.lastRunOutput, text) {
		return fmt.Errorf(`text found: %s`, text)
	}
	return nil
}

func (state *FeatureState) itPrintsTheError(expected *gherkin.DocString) error {
	if !strings.Contains(state.lastRunOutput, expected.Content) {
		return fmt.Errorf("text not found: %s\n\nactual text:\n%s", expected.Content, state.lastRunOutput)
	}
	if state.lastRunErr == nil {
		return fmt.Errorf("expected error")
	}
	return nil
}
