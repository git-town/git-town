package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"sort"
	"strings"

	"github.com/sergi/go-diff/diffmatchpatch"
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
			sortedNames := make([]string, len(constNames))
			copy(sortedNames, constNames)
			sort.Slice(sortedNames, func(i, j int) bool {
				return strings.ToLower(sortedNames[i]) < strings.ToLower(sortedNames[j])
			})

			// Use line-by-line diffmatchpatch to show complete lines
			dmp := diffmatchpatch.New()
			actualText := strings.Join(constNames, "\n")
			expectedText := strings.Join(addEmptyLinesBetweenLetters(sortedNames), "\n")

			// Convert to line-based diff
			lineArray1, lineArray2, lineHash := dmp.DiffLinesToChars(actualText, expectedText)
			diffs := dmp.DiffMain(lineArray1, lineArray2, false)
			diffs = dmp.DiffCharsToLines(diffs, lineHash)

			diffText := dmp.DiffPrettyText(diffs)

			issues = append(issues, fmt.Sprintf("%s: Constants in const block are not sorted alphabetically", filePath))
			issues = append(issues, "Differences (- actual, + expected):")
			issues = append(issues, diffText)
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

func addEmptyLinesBetweenLetters(names []string) []string {
	if len(names) == 0 {
		return names
	}

	result := []string{names[0]}
	var prevFirstLetter rune
	if len(names[0]) > 0 {
		prevFirstLetter = rune(strings.ToLower(names[0])[0])
	}

	for i := 1; i < len(names); i++ {
		if len(names[i]) > 0 {
			currentFirstLetter := rune(strings.ToLower(names[i])[0])
			if currentFirstLetter != prevFirstLetter {
				result = append(result, "")
				prevFirstLetter = currentFirstLetter
			}
		}
		result = append(result, names[i])
	}

	return result
}
