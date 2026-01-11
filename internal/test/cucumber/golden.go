package cucumber

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/git-town/git-town/v22/internal/gohacks"
	"github.com/git-town/git-town/v22/internal/gohacks/stringslice"
)

// ChangeFeatureFile updates the given section of the given feature file with the given new section.
func ChangeFeatureFile(filePath, oldSection, newSection string) error {
	// read file
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read feature file: %w", err)
	}
	fileLines := stringslice.Lines(string(content))

	// normalize file lines for searching
	normalizedFileLines := ReplaceSHAPlaceholder(fileLines)
	normalizedFileLines = ReplaceSHA(normalizedFileLines)
	normalizedFileLines = NormalizeWhitespace(normalizedFileLines)

	// normalize old section for searching
	oldSectionLines := stringslice.TrimEmptyLines(stringslice.Lines(oldSection))
	oldSectionLines = ReplaceSHAPlaceholder(oldSectionLines)
	oldSectionLines = ReplaceSHA(oldSectionLines)
	oldSectionLines = NormalizeWhitespace(oldSectionLines)

	// normalize new section
	newSectionLines := stringslice.TrimEmptyLines(stringslice.Lines(newSection))
	newSectionLines = NormalizeWhitespace(newSectionLines)

	// find the old section in the file
	startLine, found := stringslice.LocateSection(normalizedFileLines, oldSectionLines).Get()
	if !found {
		fmt.Println("WANTED SECTION START")
		fmt.Println(strings.Join(oldSectionLines, "\n"))
		fmt.Println("WANTED SECTION END")
		fmt.Println("FILE CONTENT START")
		fmt.Println(strings.Join(normalizedFileLines, "\n"))
		fmt.Println("FILE CONTENT END")
		return fmt.Errorf("could not find section in feature file %q", filePath)
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

// NormalizeWhitespace collapses redundant whitespace in the given lines.
func NormalizeWhitespace(lines []string) []string {
	return stringslice.ReplaceRegex(lines, regexp.MustCompile(`\s{2,}`), " ")
}

func ReplaceSHA(lines []string) []string {
	return stringslice.ReplaceRegex(lines, regexp.MustCompile(`[0-9a-f]{40}`), "SHA")
}

// ReplaceSHAPlaceholder replaces all placeholders like "{{ sha.* }}" with "SHA".
func ReplaceSHAPlaceholder(lines []string) []string {
	return stringslice.ReplaceRegex(lines, regexp.MustCompile(`\{\{.*?\}\}`), "SHA")
}
