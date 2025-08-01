package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"sort"
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
			sortedNames := make([]string, len(constNames))
			copy(sortedNames, constNames)
			sort.Strings(sortedNames)

			issues = append(issues, fmt.Sprintf("%s: Constants in const block are not sorted alphabetically", filePath))
			issues = append(issues, fmt.Sprintf("Expected order: %s", strings.Join(sortedNames, ", ")))
			issues = append(issues, fmt.Sprintf("Actual order: %s", strings.Join(constNames, ", ")))
		}
	}

	return issues
}

func isSorted(names []string) bool {
	for i := 1; i < len(names); i++ {
		if strings.Compare(names[i-1], names[i]) > 0 {
			return false
		}
	}
	return true
}
