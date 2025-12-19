package main

import (
	"fmt"
	"os"
	"regexp"

	"github.com/git-town/git-town/v22/pkg/asserts"
)

const (
	fileName   = "internal/test/cucumber/steps.go"
	filePath   = "../../" + fileName
	featureDir = "../../features"
)

var stepUsageRE *regexp.Regexp

func main() {
	stepsFileBytes := asserts.NoError1(os.ReadFile(filePath))
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

	unanchoredStepDefs := CheckStepRegexAnchors(existingStepDefs)
	for _, unanchoredStepDef := range unanchoredStepDefs {
		fmt.Printf("%s:%d step definition regex must start with ^ and end with $: %s\n", fileName, unanchoredStepDef.Line, unanchoredStepDef.Text)
	}

	unusedStepDefs := AllUnusedStepDefs(existingStepDefs)
	for _, unusedStepDef := range unusedStepDefs {
		fmt.Printf("%s:%d unused step definition: %s\n", fileName, unusedStepDef.Line, unusedStepDef.Text)
	}
	if len(unanchoredStepDefs) > 0 || len(unusedStepDefs) > 0 {
		os.Exit(1)
	}
}
