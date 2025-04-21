package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	gherkin "github.com/cucumber/gherkin/go/v26"
	messages "github.com/cucumber/messages/go/v21"
	"github.com/git-town/git-town/v19/pkg/asserts"
)

const (
	featureDir    = "features"                               // Subdirectory to scan
	featureExt    = ".feature"                               // File extension to look for
	targetTag     = "@messyoutput"                           // The specific tag we are interested in
	targetStepPat = `^I run ".*" and enter into the dialog$` // Regex for the step
)

// scenarioInfo holds details about a scenario relevant to our analysis.
type scenarioInfo struct {
	FilePath      string
	LineNumber    int64 // Gherkin library uses uint32 for line numbers
	HasTargetTag  bool
	HasDialogStep bool
}

func main() {
	dialogStepRegex := asserts.NoError1(regexp.Compile(targetStepPat))
	errors := 0
	asserts.NoError(filepath.WalkDir(filepath.Join(".", featureDir), func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(strings.ToLower(d.Name()), featureExt) {
			return nil
		}
		feature := parseFeatureFile(path)
		for _, scenario := range processFeature(feature, path, dialogStepRegex) {
			if scenario.HasTargetTag && !scenario.HasDialogStep {
				fmt.Printf("%s:%d  unnecessary tag\n", scenario.FilePath, scenario.LineNumber)
				errors += 1
			}
			if !scenario.HasTargetTag && scenario.HasDialogStep {
				fmt.Printf("%s:%d  missing tag\n", scenario.FilePath, scenario.LineNumber)
				errors += 1
			}
		}
		return nil
	}))
	os.Exit(errors)
}

func parseFeatureFile(filePath string) *messages.Feature {
	file := asserts.NoError1(os.Open(filePath))
	defer file.Close()
	idGenerator := messages.Incrementing{}
	document := asserts.NoError1(gherkin.ParseGherkinDocument(file, idGenerator.NewId))
	if document.Feature == nil {
		log.Fatalf("%s:  no feature definitions", filePath)
	}
	return document.Feature
}

// processFeature parses a single Gherkin file and extracts relevant scenario info.
func processFeature(feature *messages.Feature, filePath string, dialogStepRegex *regexp.Regexp) []scenarioInfo {
	result := []scenarioInfo{}
	featureHasTag := hasTag(feature.Tags, targetTag)
	backgroundHasStep := false
	for _, child := range feature.Children {
		if child.Background != nil {
			backgroundHasStep = hasStep(child.Background.Steps, dialogStepRegex)
		} else if child.Scenario != nil {
			processScenario(child.Scenario, filePath, featureHasTag, backgroundHasStep, dialogStepRegex, &result)
		} else if child.Rule != nil {
			log.Fatalf("please implement parsing the Rule's children, which are similar to the feature children")
		} else {
			fmt.Println("child has no known attributes")
		}
	}
	return result
}

// processScenario extracts information from a single scenario message.
func processScenario(scenario *messages.Scenario, filePath string, featureHasTag bool, backgroundHasStep bool, dialogStepRegex *regexp.Regexp, results *[]scenarioInfo) {
	scenarioHasTag := hasTag(scenario.Tags, targetTag)
	overallHasTag := featureHasTag || scenarioHasTag
	scenarioHasStep := hasStep(scenario.Steps, dialogStepRegex)
	overallHasStep := backgroundHasStep || scenarioHasStep
	*results = append(*results, scenarioInfo{
		FilePath:      filePath,
		LineNumber:    scenario.Location.Line,
		HasTargetTag:  overallHasTag,
		HasDialogStep: overallHasStep,
	})

	// --- Handle Scenario Outline Examples ---
	// Scenario outlines themselves don't run, their examples do.
	// The gherkin parser expands outlines into concrete scenarios implicitly
	// if you use tools like Cucumber, but the AST here represents the outline structure.
	// If the tag/step logic needs to apply to *each generated example*,
	// we'd need more complex handling, possibly involving Example tables.
	// For this request, we analyze the Scenario/Scenario Outline definition itself.
	// The line number points to the "Scenario Outline:" line.
	if len(scenario.Examples) > 0 {
		// log.Printf("ℹ️ Scenario Outline '%s' (line %d) found. Analyzing definition, not individual examples.", scenario.Name, scenario.Location.Line)
	}
}

func hasStep(steps []*messages.Step, stepRE *regexp.Regexp) bool {
	for _, step := range steps {
		if stepRE.MatchString(step.Text) {
			return true
		}
	}
	return false
}

// hasTag checks if a slice of tags contains the target tag string.
func hasTag(tags []*messages.Tag, targetTagName string) bool {
	for _, tag := range tags {
		if tag.Name == targetTagName {
			return true
		}
	}
	return false
}
