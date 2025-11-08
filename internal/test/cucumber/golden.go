package cucumber

import (
	"fmt"
	"os"
	"strings"

	"github.com/git-town/git-town/v22/internal/gohacks"
	"github.com/git-town/git-town/v22/internal/gohacks/stringslice"
)

// UpdateFeatureFile updates the given section of the given feature file with the given new section.
func UpdateFeatureFile(filePath, oldSection, newSection string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read feature file: %w", err)
	}
	fileLines := stringslice.Lines(string(content))
	oldSectionLines := stringslice.TrimEmptyLines(stringslice.Lines(oldSection))
	newSectionLines := stringslice.TrimEmptyLines(stringslice.Lines(newSection))

	// find the section in the file
	startLine, found := stringslice.LocateSection(fileLines, oldSectionLines)
	if !found {
		return fmt.Errorf("could not find section in feature file %q: %s", filePath, oldSection)
	}

	// indent the new section the same way the old one is indented in the file
	indentation := gohacks.LeadingWhitespace(fileLines[startLine])
	indentedNewSectionLines := stringslice.ChangeIndentNonEmpty(newSectionLines, indentation)

	// replace the old section with the new one
	newLines := append([]string{}, fileLines[:startLine]...)
	newLines = append(newLines, indentedNewSectionLines...)
	newLines = append(newLines, fileLines[startLine+len(oldSectionLines):]...)

	// Write back to the file
	newContent := strings.Join(newLines, "\n")
	//nolint:gosec // need permission 644 for feature files
	if err := os.WriteFile(filePath, []byte(newContent), 0o644); err != nil {
		return fmt.Errorf("failed to write feature file: %w", err)
	}
	return nil
}
