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
	oldSectionLines := stringslice.TrimEmptyLines(stringslice.Lines(oldSection))
	newSectionLines := stringslice.TrimEmptyLines(stringslice.Lines(newSection))

	// find the section in the file
	startLine, found := stringslice.LocateSection(fileLines, oldSectionLines)
	if !found {
		fmt.Println("ERROR! Could not find section in feature file: ", filePath)
		fmt.Println("Expected section:\n", oldSection)
		return
	}

	// indent the new section the same way the old one is indented in the file
	indentation := gohacks.LeadingWhitespace(fileLines[startLine])
	indentedNewSectionLines := stringslice.IndentNonEmpty(newSectionLines, indentation)

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
