package messyoutput

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
	targetStepPat = `^I r[au]n ".*" and enter into the dialogs?:$` // Regex for the step
)

type ScenarioInfo struct {
	File    string
	HasStep bool
	HasTag  bool
	Line    int64
}

func main() {
	dialogStepRegex := CompileRegex()
	errors := 0
	asserts.NoError(filepath.WalkDir(filepath.Join(".", featureDir), func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(strings.ToLower(d.Name()), featureExt) {
			return err
		}
		feature := ReadGherkinFile(path)
		for _, scenario := range FindScenarios(feature, path, dialogStepRegex) {
			if scenario.HasTag && !scenario.HasStep {
				fmt.Printf("%s:%d  unnecessary tag\n", scenario.File, scenario.Line)
				errors += 1
			}
			if !scenario.HasTag && scenario.HasStep {
				fmt.Printf("%s:%d  missing tag\n", scenario.File, scenario.Line)
				errors += 1
			}
		}
		return nil
	}))
	os.Exit(errors)
}

func CompileRegex() *regexp.Regexp {
	return regexp.MustCompile(targetStepPat)
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

func FindScenarios(feature *messages.Feature, filePath string, dialogStepRegex *regexp.Regexp) []ScenarioInfo {
	result := []ScenarioInfo{}
	featureHasTag := hasTag(feature.Tags, targetTag)
	backgroundHasStep := false
	for _, child := range feature.Children {
		switch {
		case child.Background != nil:
			backgroundHasStep = hasStep(child.Background.Steps, dialogStepRegex)
		case child.Scenario != nil:
			scenarioHasTag := hasTag(child.Scenario.Tags, targetTag)
			scenarioHasStep := hasStep(child.Scenario.Steps, dialogStepRegex)
			result = append(result, ScenarioInfo{
				File:    filePath,
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
