package main

import (
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"github.com/git-town/git-town/v14/test/asserts"
	"golang.org/x/exp/maps"
)

var unusedWhitelist = []string{ //nolint:gochecknoglobals
	`^display "([^"]+)"$`,
	`^inspect the commits$`,
	`^inspect the repo$`,
}

func FindAllUnusedStepDefs(definedSteps []StepDefinition) []StepDefinition {
	usedSteps := FindAllUsedSteps()
	return FindUnusedStepDefs(definedSteps, usedSteps)
}

func FindUnusedStepDefs(definedSteps []StepDefinition, usedSteps []string) []StepDefinition {
	unusedStepDefs := []StepDefinition{}
	for _, stepDefRE := range CreateStepRegexes(definedSteps) {
		if slices.Contains(unusedWhitelist, stepDefRE.regex.String()) {
			continue
		}
		for _, usedStep := range usedSteps {
			if stepDefRE.regex.MatchString(usedStep) {
				continue
			}
			if !IsStepDefUsed(stepDefRE, usedSteps) {
				unusedStepDefs = append(unusedStepDefs, stepDefRE.stepDef)
			}
		}
	}
	return unusedStepDefs
}

type StepRE struct {
	regex   *regexp.Regexp
	stepDef StepDefinition
}

func CreateStepRegexes(definedSteps []StepDefinition) []StepRE {
	result := make([]StepRE, len(definedSteps))

	for d, definedStep := range definedSteps {
		result[d] = StepRE{
			regex:   regexp.MustCompile(definedStep.Text),
			stepDef: definedStep,
		}
	}
	return result
}

func IsStepDefUsed(definedStep StepRE, usedSteps []string) bool {
	for _, usedStep := range usedSteps {
		if definedStep.regex.MatchString(usedStep) {
			return true
		}
	}
	return false
}

func FindAllUsedSteps() []string {
	result := map[string]bool{}
	err := filepath.WalkDir(featureDir, func(path string, _ os.DirEntry, err error) error {
		asserts.NoError(err)
		if filepath.Ext(path) != ".feature" {
			return nil
		}
		fileContent, err := os.ReadFile(path)
		asserts.NoError(err)
		for _, stepInFile := range FindUsedStepsIn(string(fileContent)) {
			result[stepInFile] = true
		}
		return nil
	})
	asserts.NoError(err)
	return maps.Keys(result)
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
