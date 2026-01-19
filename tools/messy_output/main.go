package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"

	gherkin "github.com/cucumber/gherkin/go/v26"
	messages "github.com/cucumber/messages/go/v21"
	"github.com/git-town/git-town/v22/pkg/asserts"
)

const tag = "@messyoutput"

type ScenarioInfo struct {
	HasStep bool
	HasTag  bool
	Line    int64
}

func main() {
	dialogStepRegex := regexp.MustCompile(`^I r[au]n ".*" and enter into the dialogs?:$`)
	errors := 0
	asserts.NoError(filepath.WalkDir("features", func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() || filepath.Ext(d.Name()) != ".feature" {
			return err
		}
		feature := ReadGherkinFile(path)
		scenarios := FindScenarios(feature, dialogStepRegex)
		for _, error := range AnalyzeScenarios(path, scenarios) {
			fmt.Println(error)
			errors += 1
		}
		return nil
	}))
	os.Exit(errors)
}

func AnalyzeScenarios(file string, scenarios []ScenarioInfo) []string {
	result := []string{}
	for _, scenario := range scenarios {
		if scenario.HasTag && !scenario.HasStep {
			result = append(result, fmt.Sprintf("%s:%d  unnecessary @messyoutput tag\n", file, scenario.Line))
		}
		if !scenario.HasTag && scenario.HasStep {
			result = append(result, fmt.Sprintf("%s:%d  missing @messyoutput tag\n", file, scenario.Line))
		}
	}
	return result
}

func FindScenarios(feature *messages.Feature, dialogStepRegex *regexp.Regexp) []ScenarioInfo {
	result := []ScenarioInfo{}
	featureHasTag := hasTag(feature.Tags, tag)
	backgroundHasStep := false
	for _, child := range feature.Children {
		switch {
		case child.Background != nil:
			backgroundHasStep = hasStep(child.Background.Steps, dialogStepRegex)
		case child.Scenario != nil:
			scenarioHasTag := hasTag(child.Scenario.Tags, tag)
			scenarioHasStep := hasStep(child.Scenario.Steps, dialogStepRegex)
			result = append(result, ScenarioInfo{
				HasStep: backgroundHasStep || scenarioHasStep,
				HasTag:  featureHasTag || scenarioHasTag,
				Line:    child.Scenario.Location.Line,
			})
		case child.Rule != nil:
			log.Fatalf("please implement parsing the Rule's children, which are similar to the feature children")
		default:
			fmt.Println("child has no known attributes")
		}
	}
	return result
}

// parses the content of a Gherkin file
func ReadGherkinFile(filePath string) *messages.Feature {
	file := asserts.NoError1(os.Open(filePath))
	defer file.Close()
	idGenerator := messages.Incrementing{}
	document := asserts.NoError1(gherkin.ParseGherkinDocument(file, idGenerator.NewId))
	if document.Feature == nil {
		log.Fatalf("%s:  no feature definitions", filePath) //nolint:gocritic
	}
	return document.Feature
}

// indicates whether the given steps contain a step that matches the given regex
func hasStep(steps []*messages.Step, stepRE *regexp.Regexp) bool {
	for _, step := range steps {
		if stepRE.MatchString(step.Text) {
			return true
		}
	}
	return false
}

// indicates whether the given tags contain a tag with the given name
func hasTag(tags []*messages.Tag, name string) bool {
	for _, tag := range tags {
		if tag.Name == name {
			return true
		}
	}
	return false
}
