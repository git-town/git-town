package steps

import (
	"fmt"
	"strings"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/gherkin"
)

// PrintSteps defines Gherkin steps around printing things to the terminal.
func PrintSteps(suite *godog.Suite, fs *FeatureState) {
	suite.Step(`^it does not print "([^\"]*)"$`, func(text string) error {
		if strings.Contains(fs.activeScenarioState.lastRunResult.OutputSanitized(), text) {
			return fmt.Errorf("text found: %q", text)
		}
		return nil
	})

	suite.Step(`^it prints$`, func(expected *gherkin.DocString) error {
		if !strings.Contains(fs.activeScenarioState.lastRunResult.OutputSanitized(), expected.Content) {
			return fmt.Errorf("text not found:\n\nEXPECTED: %q\n\nACTUAL:\n\n%q", expected.Content, fs.activeScenarioState.lastRunResult.OutputSanitized())
		}
		return nil
	})

	suite.Step(`^it prints the error:$`, func(expected *gherkin.DocString) error {
		if !strings.Contains(fs.activeScenarioState.lastRunResult.OutputSanitized(), expected.Content) {
			return fmt.Errorf("text not found: %s\n\nactual text:\n%s", expected.Content, fs.activeScenarioState.lastRunResult.OutputSanitized())
		}
		if fs.activeScenarioState.lastRunErr == nil {
			return fmt.Errorf("expected error")
		}
		return nil
	})

	suite.Step(`^it prints the initial configuration prompt$`, func() error {
		expected := "Git Town needs to be configured"
		if !fs.activeScenarioState.lastRunResult.OutputContainsText(expected) {
			return fmt.Errorf("text not found:\n\nEXPECTED: %q\n\nACTUAL:\n\n%q\n----------------------------", expected, fs.activeScenarioState.lastRunResult.Output())
		}
		return nil
	})
}
