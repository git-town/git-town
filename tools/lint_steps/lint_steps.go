package main

import (
	"fmt"
	"os"
	"regexp"

	"github.com/git-town/git-town/v14/test/asserts"
)

const (
	fileName   = "test/cucumber/steps.go"
	filePath   = "../../" + fileName
	featureDir = "../../features"
)

var stepUsageRE *regexp.Regexp //nolint:gochecknoglobals

func main() {
	stepsFileBytes, err := os.ReadFile(filePath)
	asserts.NoError(err)
	stepsFileText := string(stepsFileBytes)

	malformattedStepDefs := CheckStepDefinitions(stepsFileText)
	for _, issue := range malformattedStepDefs {
		fmt.Printf("%s:%d step definition must use backticks\n", fileName, issue.Line)
	}
	if len(malformattedStepDefs) > 0 {
		os.Exit(1)
	}

	existingStepDefs := FindStepDefinitions(stepsFileText)
	if len(existingStepDefs) == 0 {
		panic("no step definitions found")
	}

	unsortedStepDefs := AllUnsortedStepDefs(existingStepDefs)
	for _, unsortedStepDef := range unsortedStepDefs {
		fmt.Printf("%s:%d steps are not alphabetically sorted, expected here: %s\n", fileName, unsortedStepDef.Line, unsortedStepDef.Text)
	}

	unusedStepDefs := AllUnusedStepDefs(existingStepDefs)
	for _, unusedStepDef := range unusedStepDefs {
		fmt.Printf("%s:%d unused step definition: %s\n", fileName, unusedStepDef.Line, unusedStepDef.Text)
	}
	if len(unsortedStepDefs) > 0 || len(unusedStepDefs) > 0 {
		os.Exit(1)
	}
}
