package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"slices"
	"sort"
	"strings"

	"github.com/cucumber/godog/colors"
)

func main() {
	issues := findUnsortedConst("internal/messages/en.go")
	if len(issues) > 0 {
		for _, issue := range issues {
			fmt.Println(issue)
		}
		os.Exit(1)
	}
}

func findUnsortedConst(filePath string) []string {
	var issues []string

	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, filePath, nil, parser.ParseComments)
	if err != nil {
		issues = append(issues, fmt.Sprintf("Error parsing %s: %v", filePath, err))
		return issues
	}

	// Find the const block
	for _, decl := range file.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.CONST {
			continue
		}

		// Extract constant names from the const block
		var constNames []string
		for _, spec := range genDecl.Specs {
			valueSpec, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}
			for _, name := range valueSpec.Names {
				constNames = append(constNames, name.Name)
			}
		}

		// Check if constants are sorted alphabetically
		if !isSorted(constNames) {
			sortedNames := make([]string, len(constNames))
			copy(sortedNames, constNames)
			sort.Slice(sortedNames, func(i, j int) bool {
				return strings.ToLower(sortedNames[i]) < strings.ToLower(sortedNames[j])
			})

			unsortedPositions := findUnsortedPositions(constNames)
			diffLines := getDiffContextLines(filePath, constNames, sortedNames, unsortedPositions, 3)

			issues = append(issues, "Constants are not sorted alphabetically")
			for _, line := range diffLines {
				issues = append(issues, line)
			}
		}
	}

	return issues
}

func findUnsortedPositions(names []string) []int {
	var unsortedPositions []int

	for i := 1; i < len(names); i++ {
		if strings.Compare(strings.ToLower(names[i-1]), strings.ToLower(names[i])) > 0 {
			// Both the previous and current items are part of the unsorted sequence
			if len(unsortedPositions) == 0 || unsortedPositions[len(unsortedPositions)-1] != i-1 {
				unsortedPositions = append(unsortedPositions, i-1)
			}
			unsortedPositions = append(unsortedPositions, i)
		}
	}

	return unsortedPositions
}

func getDiffContextLines(filePath string, actualNames []string, expectedNames []string, unsortedPositions []int, contextSize int) []string {
	if len(unsortedPositions) == 0 {
		return []string{}
	}

	var result []string
	covered := make(map[int]bool)

	// For each unsorted position, show it with context
	for _, pos := range unsortedPositions {
		if covered[pos] {
			continue
		}

		// Calculate context range
		start := pos - contextSize
		start = max(start, 0)
		end := pos + contextSize
		if end >= len(actualNames) {
			end = len(actualNames) - 1
		}

		// Add separator if this isn't the first group
		if len(result) > 0 {
			result = append(result, "--")
		}

		// Add lines with context, showing diff format
		for i := start; i <= end; i++ {
			covered[i] = true
			if slices.Contains(unsortedPositions, i) {
				// Show the diff for unsorted lines
				result = append(result, colors.Red(fmt.Sprintf("%s:%d: %s", filePath, i+4, actualNames[i])))
				result = append(result, colors.Green(fmt.Sprintf("%s:%d: %s", filePath, i+4, expectedNames[i])))
			} else {
				// Show context lines without prefix
				result = append(result, fmt.Sprintf("%s:%d: %s", filePath, i+4, actualNames[i]))
			}
		}
	}

	return result
}

func isSorted(names []string) bool {
	for i := 1; i < len(names); i++ {
		if strings.Compare(strings.ToLower(names[i-1]), strings.ToLower(names[i])) > 0 {
			return false
		}
	}
	return true
}
