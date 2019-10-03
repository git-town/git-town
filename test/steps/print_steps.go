package steps

import (
	"fmt"
	"strings"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
)

// PrintSteps defines Gherkin steps around printing things to the terminal.
func PrintSteps(s *godog.Suite, gtf *GitTownFeature) {
	s.Step("^it prints$", gtf.itPrints)
	s.Step("^it does not print \"([^\"]*)\"$", gtf.itDoesNotPrint)
	s.Step(`^it prints the error:$`, gtf.itPrintsTheError)
}

func (gtf *GitTownFeature) itPrints(expected *gherkin.DocString) error {
	if !strings.Contains(gtf.lastRunOutput, expected.Content) {
		return fmt.Errorf(`text not found: %s`, expected.Content)
	}
	return nil
}

func (gtf *GitTownFeature) itDoesNotPrint(text string) error {
	if strings.Contains(gtf.lastRunOutput, text) {
		return fmt.Errorf(`text found: %s`, text)
	}
	return nil
}

func (gtf *GitTownFeature) itPrintsTheError(expected *gherkin.DocString) error {
	if !strings.Contains(gtf.lastRunOutput, expected.Content) {
		return fmt.Errorf("text not found: %s\n\nactual text:\n%s", expected.Content, gtf.lastRunOutput)
	}
	if gtf.lastRunErr == nil {
		return fmt.Errorf("expected error")
	}
	return nil
}
