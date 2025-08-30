package main

import (
	"regexp"
	"strings"
)

func CheckStepDefinitions(fileContent string) []StepDefinition {
	var result []StepDefinition
	re := regexp.MustCompile(`sc\.Step\(['"]`)
	for l, line := range strings.Split(fileContent, "\n") {
		match := re.FindAllString(line, -1)
		if len(match) > 0 {
			result = append(result, StepDefinition{
				Line: l + 1,
				Text: match[0],
			})
		}
	}
	return result
}
