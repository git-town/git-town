package main

import "strings"

// CheckStepRegexAnchors verifies that all step definition regular expressions
// start with "^" and end with "$" to ensure proper matching.
func CheckStepRegexAnchors(stepDefs []StepDefinition) []StepDefinition {
	var result []StepDefinition
	for _, stepDef := range stepDefs {
		if !strings.HasPrefix(stepDef.Text, "^") || !strings.HasSuffix(stepDef.Text, "$") {
			result = append(result, stepDef)
		}
	}
	return result
}
