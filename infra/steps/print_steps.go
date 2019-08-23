package steps

import (
	"fmt"
	"strings"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
)

func PrintSteps(s *godog.Suite) {
	s.Step("^it prints$", func(expected *gherkin.DocString) error {
		if !strings.Contains(lastRunResult.Output, expected.Content) {
			return fmt.Errorf(`text not found: %s`, expected.Content)
		}
		return nil
	})

	s.Step("^it does not print \"([^\"]*)\"$",
		func(text string) error {
			if strings.Contains(lastRunResult.Output, text) {
				return fmt.Errorf(`text found: %s`, text)
			}
			return nil
		})

	s.Step(`^it prints the error:$`, func(expected *gherkin.DocString) error {
		if !strings.Contains(lastRunResult.Output, expected.Content) {
			return fmt.Errorf("text not found: %s\n\nactual text:\n%s", expected.Content, lastRunResult.Output)
		}
		if lastRunResult.Err == nil {
			return fmt.Errorf("expected error")
		}
		return nil
	})
}
