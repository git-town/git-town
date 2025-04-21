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
	ScenarioName  string
	HasTargetTag  bool
	HasDialogStep bool
}

func main() {
	dialogStepRegex := asserts.NoError1(regexp.Compile(targetStepPat))
	scenarios := []scenarioInfo{}
	targetPath := filepath.Join(".", featureDir) // Look in ./feature
	asserts.NoError(filepath.WalkDir(targetPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			log.Printf("cannot access path %q: %v\n", path, err)
			return err
		}
		if d.IsDir() || !strings.HasSuffix(strings.ToLower(d.Name()), featureExt) {
			return nil
		}
		// Process the found feature file
		scenario := asserts.NoError1(processFeatureFile(path, dialogStepRegex))
		scenarios = append(scenarios, scenario...)
		return nil
	}))
	for _, scenario := range scenarios {
		if scenario.HasTargetTag && !scenario.HasDialogStep {
			fmt.Printf("%s:%d  unnecessary tag\n", scenario.FilePath, scenario.LineNumber)
		}
		if !scenario.HasTargetTag && scenario.HasDialogStep {
			fmt.Printf("%s:%d  missing tag\n", scenario.FilePath, scenario.LineNumber)
		}
	}
	os.Exit(len(scenarios))
}

// processFeatureFile parses a single Gherkin file and extracts relevant scenario info.
func processFeatureFile(filePath string, dialogStepRegex *regexp.Regexp) ([]scenarioInfo, error) {
	file := asserts.NoError1(os.Open(filePath))
	defer file.Close()
	result := []scenarioInfo{}
	idGenerator := messages.Incrementing{}
	document := asserts.NoError1(gherkin.ParseGherkinDocument(file, idGenerator.NewId))
	if document.Feature == nil {
		log.Printf("ℹ️ File %s does not contain a Feature definition. Skipping.", filePath)
		return result, nil // Not an error, just an empty file or non-feature content
	}
	featureHasTag := hasTag(document.Feature.Tags, targetTag)
	// --- Iterate through Scenarios and Scenario Outlines ---
	for _, child := range document.Feature.Children {
		var scenario *messages.Scenario
		if child.Scenario != nil {
			scenario = child.Scenario
		} else if child.Rule != nil {
			// Also check scenarios within Rules
			rule := child.Rule
			for _, ruleChild := range rule.Children {
				if ruleChild.Scenario != nil {
					scenario = ruleChild.Scenario
					// Process the scenario found within the rule
					processSingleScenario(scenario, filePath, featureHasTag, dialogStepRegex, &result)
				}
			}
			// Skip to next top-level child after processing rule's scenarios
			continue
		}
		if scenario != nil {
			processSingleScenario(scenario, filePath, featureHasTag, dialogStepRegex, &result)
		}
	}
	return result, nil
}

// processSingleScenario extracts information from a single scenario message.
func processSingleScenario(scenario *messages.Scenario, filePath string, featureHasTargetTag bool, dialogStepRegex *regexp.Regexp, results *[]scenarioInfo) {
	// Check scenario-level tag
	scenarioHasTargetTag := hasTag(scenario.Tags, targetTag)

	// Determine overall tag status (feature OR scenario)
	overallHasTargetTag := featureHasTargetTag || scenarioHasTargetTag

	// --- Check for Target Step ---
	scenarioHasDialogStep := false
	for _, step := range scenario.Steps {
		if dialogStepRegex.MatchString(step.Text) {
			scenarioHasDialogStep = true
			break // Found the step, no need to check further steps in this scenario
		}
	}

	// --- Store Scenario Info ---
	info := scenarioInfo{
		FilePath:      filePath,
		ScenarioName:  scenario.Name,
		LineNumber:    scenario.Location.Line,
		HasTargetTag:  overallHasTargetTag,
		HasDialogStep: scenarioHasDialogStep,
	}
	*results = append(*results, info)

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

// hasTag checks if a slice of tags contains the target tag string.
func hasTag(tags []*messages.Tag, targetTagName string) bool {
	for _, tag := range tags {
		if tag.Name == targetTagName {
			return true
		}
	}
	return false
}
