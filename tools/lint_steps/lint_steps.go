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
	stepDefs := FindStepDefinitions(fileContent)
	for s, stepDef := range stepDefs {
		fmt.Printf("%d. %s\n", s+1, stepDef)
	}
}

func FindStepDefinitions(fileContent string) []string {
	result := []string{}
	re := regexp.MustCompile("^\\s*sc\\.Step\\(`(.*)`")
	matches := re.FindAllStringSubmatch(fileContent, -1)
	for _, match := range matches {
		if len(match) > 1 {
			result = append(result, match[1])
		}
	}
	return result
}
