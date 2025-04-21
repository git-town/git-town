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
	featureDir    = "features"                                     // Subdirectory to scan
	featureExt    = ".feature"                                     // File extension to look for
	targetTag     = "@messyoutput"                                 // The specific tag we are interested in
	targetStepPat = `^I r[ua]n ".*" and enter into the dialogs?:$` // Regex for the step
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
		// path := "./features/append/on_feature_branch/missing_lineage/unknown_parent.feature"
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
	// fmt.Println("background has step:", backgroundHasStep)
	// fmt.Println("scenario has step:", scenarioHasStep)
	// fmt.Println("feature has tag:", featureHasTag)
	// fmt.Println("overall has tag:", overallHasTag)
	*results = append(*results, scenarioInfo{
		FilePath:      filePath,
		LineNumber:    scenario.Location.Line,
		HasTargetTag:  overallHasTag,
		HasDialogStep: overallHasStep,
	})
	if len(scenario.Examples) > 0 {
		// log.Printf("ℹ️ Scenario Outline '%s' (line %d) found. Analyzing definition, not individual examples.", scenario.Name, scenario.Location.Line)
	}
}

func hasStep(steps []*messages.Step, stepRE *regexp.Regexp) bool {
	for _, step := range steps {
		// fmt.Println("step:", step.Text)
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
