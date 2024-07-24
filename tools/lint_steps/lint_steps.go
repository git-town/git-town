package main

import (
	"fmt"
	"os"
	"regexp"

	"github.com/git-town/git-town/v14/test/asserts"
)

func main() {
	content, err := os.ReadFile("test/cucumber/steps.go")
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
	for s, stepDef := range stepDefs {
		fmt.Printf("%d. %s\n", s+1, stepDef)
	}
}

func CheckStepDefinitions(fileContent string) []string {
	re := regexp.MustCompile(`sc\.Step\(['"]`)
	return re.FindAllString(fileContent, -1)
}

func FindStepDefinitions(fileContent string) []string {
	result := []string{}
	re := regexp.MustCompile("sc\\.Step\\(`(.*)`")
	matches := re.FindAllStringSubmatch(fileContent, -1)
	for _, match := range matches {
		result = append(result, match[1])
	}
	return result
}
