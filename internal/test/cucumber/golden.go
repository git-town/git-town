package cucumber

import (
	"fmt"
	"os"
	"strings"

	"github.com/git-town/git-town/v22/internal/gohacks"
	"github.com/git-town/git-town/v22/internal/gohacks/stringslice"
)

// UpdateFeatureFile updates the given section of the given feature file with the given new section.
func UpdateFeatureFile(filePath, oldSection, newSection string) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		panic(fmt.Sprintf("failed to read feature file: %v", err))
	}
	fileLines := strings.Split(string(content), "\n")
	oldSectionLines := trimmedLines(oldSection)
	newSectionLines := trimmedLines(newSection)

	// find the section in the file
	startLine, found := stringslice.LocateSection(fileLines, oldSectionLines)
	if !found {
		fmt.Println("ERROR! Could not find section in feature file: ", filePath)
		fmt.Println("Expected section:\n", oldSection)
		return
	}

	// indent the new section the same way the old one is indented in the file
	indentation := gohacks.LeadingWhitespace(fileLines[startLine])
	indentedNewSectionLines := IndentSection(newSectionLines, indentation)

	// replace the old section with the new one
	newLines := append([]string{}, fileLines[:startLine]...)
	newLines = append(newLines, indentedNewSectionLines...)
	newLines = append(newLines, fileLines[startLine+len(oldSectionLines):]...)

	// Write back to the file
	newContent := strings.Join(newLines, "\n")
	//nolint:gosec // need permission 644 for feature files
	if err := os.WriteFile(filePath, []byte(newContent), 0o644); err != nil {
		panic(fmt.Sprintf("failed to write feature file: %v", err))
	}
}

// IndentSection applies indentation to each non-empty line
func IndentSection(lines []string, indentation string) []string {
	result := make([]string, len(lines))
	for i, line := range lines {
		if strings.TrimSpace(line) != "" {
			result[i] = indentation + strings.TrimLeft(line, " \t")
		} else {
			result[i] = line
		}
	}
	return result
}

// trimmedLines removes leading and trailing empty lines from a table string
func trimmedLines(tableStr string) []string {
	linesRaw := strings.Split(tableStr, "\n")
	return stringslice.TrimEmptyLines(linesRaw)
}
