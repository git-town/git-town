package main

import (
	"fmt"
	"os"
	"regexp"

	"github.com/git-town/git-town/v14/test/asserts"
)

const fileName = "test/cucumber/steps.go"
const filePath = "../../" + fileName
const featureDir = "../../features"

var stepUsageRE *regexp.Regexp

func main() {
	content, err := os.ReadFile(filePath)
	asserts.NoError(err)
	fileContent := string(content)

	if malformattedStepDefs := CheckStepDefinitions(fileContent); len(malformattedStepDefs) > 0 {
		for i, issue := range malformattedStepDefs {
			fmt.Printf("%d. %s", i, issue)
		}
		os.Exit(1)
	}

	existingStepDefs := FindStepDefinitions(fileContent)
	if len(existingStepDefs) == 0 {
		panic("no step definitions found")
	}

	unsortedStepDefs := FindUnsortedStepDefs(existingStepDefs)
	if len(unsortedStepDefs) > 0 {
		for _, unsortedStepDef := range unsortedStepDefs {
			fmt.Printf("%s:%d unsorted, expected %q here\n", fileName, unsortedStepDef.Line, unsortedStepDef.Text)
		}
		os.Exit(1)
	}

	unusedStepDefs := FindAllUnusedStepDefs()
	if len(unusedStepDefs) > 0 {
		for _, unusedStepDef := range unusedStepDefs {
			fmt.Printf("%s:%d unused step definition: %s\n", fileName, unusedStepDef.Line, unusedStepDef.Text)
		}
		os.Exit(1)
	}
}
