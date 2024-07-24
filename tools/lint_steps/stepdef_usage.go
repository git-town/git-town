package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	sets "github.com/deckarep/golang-set/v2"
	"github.com/git-town/git-town/v14/test/asserts"
)

func FindAllUnusedStepDefs() []string {
	result := []string{}
	definedSteps := FindStepDefinitions(fileName)
	definedREs := make([]*regexp.Regexp, len(definedSteps))
	for d, definedStep := range definedSteps {
		definedREs[d] = regexp.MustCompile(definedStep.Text)
	}
	usedSteps := FindAllUsedSteps()
REs:
	for _, definedRE := range definedREs {
		for _, usedStep := range usedSteps {
			if definedRE.MatchString(usedStep) {
				continue REs
			}
		}
		fmt.Printf("unused step definition: %s", definedRE)
	}
	return result
}

func FindUnusedStepDefs(definedSteps []StepDefinition, usedSteps []string) []StepDefinition {
	result := []StepDefinition{}
	definedREs := make([]StepRE, len(definedSteps))
	for d, definedStep := range definedSteps {
		definedREs[d] = StepRE{
			regex:   regexp.MustCompile(definedStep.Text),
			stepDef: definedStep,
		}
	}
REs:
	for _, definedRE := range definedREs {
		for _, usedStep := range usedSteps {
			if definedRE.regex.MatchString(usedStep) {
				continue REs
			}
		}
		result = append(result, definedRE.stepDef)
	}
	return result
}

type StepRE struct {
	regex   *regexp.Regexp
	stepDef StepDefinition
}

func FindAllUsedSteps() []string {
	result := sets.NewSet[string]()
	err := filepath.WalkDir(featureDir, func(path string, entry os.DirEntry, err error) error {
		asserts.NoError(err)
		if filepath.Ext(path) != ".feature" {
			return nil
		}
		fileContent, err := os.ReadFile(path)
		asserts.NoError(err)
		stepsInFile := FindUsedStepsIn(string(fileContent))
		result.Append(stepsInFile...)
		return nil
	})
	asserts.NoError(err)
	return result.ToSlice()
}

// provides all usages of Cucumber steps in the given file content
func FindUsedStepsIn(fileContent string) []string {
	if stepUsageRE == nil {
		stepUsageRE = regexp.MustCompile(`^\s*(?:Given|When|Then|And) (.*)$`)
	}
	result := []string{}
	for _, line := range strings.Split(fileContent, "\n") {
		matches := stepUsageRE.FindAllStringSubmatch(line, -1)
		if len(matches) > 0 {
			result = append(result, strings.TrimSpace(matches[0][1]))
		}
	}
	return result
}
