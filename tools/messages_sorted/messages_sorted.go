package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
)

func main() {
	issues := lintMessagesFile("internal/messages/en.go")
	if len(issues) > 0 {
		for _, issue := range issues {
			fmt.Println(issue)
		}
		os.Exit(1)
	}
}

func lintMessagesFile(filePath string) []string {
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
			unsortedPositions := findUnsortedPositions(constNames)
			contextLines := getContextLines(constNames, unsortedPositions, 3)

			issues = append(issues, fmt.Sprintf("%s: Constants in const block are not sorted alphabetically", filePath))
			issues = append(issues, "Unsorted constants with context:")
			for _, line := range contextLines {
				issues = append(issues, line)
			}
		}
	}

	return issues
}

func isSorted(names []string) bool {
	for i := 1; i < len(names); i++ {
		if strings.Compare(strings.ToLower(names[i-1]), strings.ToLower(names[i])) > 0 {
			return false
		}
	}
	return true
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

func getContextLines(names []string, unsortedPositions []int, contextSize int) []string {
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
		if start < 0 {
			start = 0
		}
		end := pos + contextSize
		if end >= len(names) {
			end = len(names) - 1
		}
		
		// Add separator if this isn't the first group
		if len(result) > 0 {
			result = append(result, "--")
		}
		
		// Add lines with context, marking unsorted ones
		for i := start; i <= end; i++ {
			covered[i] = true
			prefix := "  "
			if contains(unsortedPositions, i) {
				prefix = "> " // Mark unsorted lines
			}
			result = append(result, fmt.Sprintf("%s%d: %s", prefix, i+1, names[i]))
		}
	}
	
	return result
}

func contains(slice []int, item int) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}
