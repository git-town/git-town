package main

import (
	"cmp"
	"regexp"
	"slices"
	"strings"
)

type StepDefinition struct {
	Line int
	Text string
}

func AllUnsortedStepDefs(stepDefs []StepDefinition) []StepDefinition {
	result := []StepDefinition{}
	sortedStepDefs := make([]string, len(stepDefs))
	for s, stepDef := range stepDefs {
		sortedStepDefs[s] = stepDef.Text
	}
	slices.SortFunc(sortedStepDefs, normalizedSort)
	for s := range sortedStepDefs {
		if stepDefs[s].Text != sortedStepDefs[s] {
			result = append(result, StepDefinition{
				Line: stepDefs[s].Line,
				Text: sortedStepDefs[s],
			})
		}
	}
	return result
}

func FindStepDefinitions(fileContent string) []StepDefinition {
	result := []StepDefinition{}
	re := regexp.MustCompile("sc\\.Step\\(`(.*)`")
	for l, line := range strings.Split(fileContent, "\n") {
		for _, match := range re.FindAllStringSubmatch(line, -1) {
			result = append(result, StepDefinition{
				Line: l + 1,
				Text: match[1],
			})
		}
	}
	return result
}

func NormalizeForSort(text string) string {
	text = strings.ToLower(text)
	for _, c := range "\"()[]^$*+?: " {
		text = strings.ReplaceAll(text, string(c), "")
	}
	return text
}

func normalizedSort(a, b string) int {
	return cmp.Compare(NormalizeForSort(a), NormalizeForSort(b))
}
