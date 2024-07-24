package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"github.com/git-town/git-town/v14/test/asserts"
)

const fileName = "test/cucumber/steps.go"
const featureDir = "features"

var stepUsageRE *regexp.Regexp

func main() {
	content, err := os.ReadFile(fileName)
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
			fmt.Printf("%s:%d expected %q", fileName, unsortedStepDef.Line, unsortedStepDef.Text)
		}
	}

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

func FindAllUsedSteps() []string {
	result := []string{}
	err := filepath.WalkDir(featureDir, func(path string, entry os.DirEntry, err error) error {
		asserts.NoError(err)
		if filepath.Ext(path) != ".feature" {
			return nil
		}
		fileContent, err := os.ReadFile(path)
		asserts.NoError(err)
		stepsInFile := FindUsedStepsIn(string(fileContent))
		result = append(result, stepsInFile...)
		return nil
	})
	asserts.NoError(err)
	return result
}

// provides all usages of Cucumber steps in the given file content
func FindUsedStepsIn(fileContent string) []string {
	initializeRE()
	result := []string{}
	for _, line := range strings.Split(fileContent, "\n") {
		matches := stepUsageRE.FindAllStringSubmatch(line, -1)
		if len(matches) > 0 {
			result = append(result, strings.TrimSpace(matches[0][1]))
		}
	}
	return result
}

func initializeRE() {
	if stepUsageRE == nil {
		stepUsageRE = regexp.MustCompile(`^\s*(?:Given|When|Then|And) (.*)$`)
	}
}
