package main

import "regexp"

func CheckStepDefinitions(fileContent string) []string {
	re := regexp.MustCompile(`sc\.Step\(['"]`)
	return re.FindAllString(fileContent, -1)
}
