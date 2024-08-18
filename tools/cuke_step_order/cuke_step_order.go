package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"

	. "github.com/git-town/git-town/v15/pkg/prelude"
	"github.com/git-town/git-town/v15/test/asserts"
)

const defs_file = "expected_steps.txt"

var stepPrefixes []string = []string{"Given", "When", "Then", "And"}

func main() {
	wantIndex := indexSteps(readFileLines(defs_file))
	featureFiles := findFeatureFiles("../../features")
	fmt.Printf("checking %v feature files\n", len(featureFiles))
	for _, featureFile := range featureFiles {
		lines := readFileLines(featureFile)
		scenarios := getScenarios(lines)
		fmt.Printf("file %q has %d scenarios\n", featureFile, len(scenarios))
		for _, scenario := range scenarios {
			checkScenarioSteps(featureFile, scenario, wantIndex, scenario)
		}
	}
}

func findFeatureFiles(dir string) []string {
	var result []string
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		asserts.NoError(err)
		if filepath.Ext(path) == ".feature" {
			result = append(result, path)
		}
		return nil
	})
	asserts.NoError(err)
	return result
}

func readFileLines(filename string) []string {
	bytes, err := os.ReadFile(defs_file)
	if err != nil {
		log.Fatalf("cannot read file %q: %v", filename, err)
	}
	return strings.Split(string(bytes), "\n")
}

func indexSteps(steps []string) map[string]int {
	result := map[string]int{}
	for s, step := range steps {
		result[step] = s
	}
	return result
}

func getScenarios(lines []string) [][]string {
	result := [][]string{}
	currentScenarioOpt := None[[]string]()
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if IsScenarioLine(line) {
			if currentScenario, hasCurrentScenario := currentScenarioOpt.Get(); hasCurrentScenario {
				result = append(result, currentScenario)
			}
		}
	}
	return result
}

func IsScenarioLine(line string) bool {
	if strings.HasPrefix(line, "Scenario: ") {
		return true
	}
	return false
}

func isStepLine(line string) bool {
	for _, prefix := range stepPrefixes {
		if strings.HasPrefix(line, prefix) {
			return true
		}
	}
	return false
}

func checkScenarioSteps(filename string, steps []string, wantIndex map[string]int, wantSteps []string) {
	haveIndex := make([]int, len(steps))
	for s, step := range steps {
		stepIndex, hasStep := wantIndex[step]
		if !hasStep {
			log.Fatalf("step not listed in order file: %s", step)
		}
		haveIndex[s] = stepIndex
	}
	if slices.IsSorted(haveIndex) {
		return
	}
	fmt.Printf("%s is unsorted, expected step order:\n", filename)
	sorted := make([]int, len(steps))
	copy(sorted, haveIndex)
	slices.Sort(sorted)
	for s := range sorted {
		fmt.Println(wantSteps[s])
	}
}
