package main

import (
	"fmt"
	"os"
	"regexp"
	"slices"
	"strings"

	"github.com/git-town/git-town/v14/test/asserts"
)

const fileName = "test/cucumber/steps.go"

func main() {
	content, err := os.ReadFile(fileName)
	asserts.NoError(err)
	fileContent := string(content)

	// ensure step definitions use backticks
	if malformattedStepDefs := CheckStepDefinitions(fileContent); len(malformattedStepDefs) > 0 {
		for i, issue := range malformattedStepDefs {
			fmt.Printf("%d. %s", i, issue)
		}
		os.Exit(1)
	}

	// find step definitions
	stepDefs := FindStepDefinitions(fileContent)
	if len(stepDefs) == 0 {
		panic("no step definitions found")
	}

	// find unsorted step definitions
	FindUnsortedStepDefs(stepDefs)

	// find unused step definitions

	// for s, stepDef := range stepDefs {
	// 	fmt.Printf("%d. %s\n", s+1, stepDef)
	// }
}

func CheckStepDefinitions(fileContent string) []string {
	re := regexp.MustCompile(`sc\.Step\(['"]`)
	return re.FindAllString(fileContent, -1)
}

func FindStepDefinitions(fileContent string) []StepDefinition {
	result := []StepDefinition{}
	re := regexp.MustCompile("sc\\.Step\\(`(.*)`")
	lines := strings.Split(fileContent, "\n")
	for l, line := range lines {
		matches := re.FindAllStringSubmatch(line, -1)
		for _, match := range matches {
			result = append(result, StepDefinition{
				Text: match[1],
				Line: l,
			})
		}
	}
	return result
}

type StepDefinition struct {
	Text string
	Line int
}

func FindUnsortedStepDefs(stepDefs []StepDefinition) []StepDefinition {
	result := []StepDefinition{}
	sortedStepDefs := make([]string, len(stepDefs))
	for s, stepDef := range stepDefs {
		sortedStepDefs[s] = stepDef.Text
	}
	slices.Sort(sortedStepDefs)
	for s := range sortedStepDefs {
		if stepDefs[s].Text != sortedStepDefs[s] {
			result = append(result, StepDefinition{
				Text: sortedStepDefs[s],
				Line: stepDefs[s].Line,
			})
		}
	}
	return result
}
