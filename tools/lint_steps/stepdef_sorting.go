package main

import (
	"cmp"
	"regexp"
	"slices"
	"strings"
)

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
	slices.SortFunc(sortedStepDefs, func(a, b string) int { return cmp.Compare(normalizeForSort(a), normalizeForSort(b)) })
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

func normalizeForSort(text string) string {
	text = strings.ToLower(text)
	for _, c := range "\"()[]^$*+ " {
		text = strings.ReplaceAll(text, string(c), "")
	}
	return text
}
