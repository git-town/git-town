package main

import (
	"regexp"
	"strings"
)

type StepDefinition struct {
	Line int
	Text string
}

func FindStepDefinitions(fileContent string) []StepDefinition {
	var result []StepDefinition
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
