package cucumber

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/git-town/git-town/v22/internal/gohacks"
	"github.com/git-town/git-town/v22/internal/gohacks/stringslice"
)

// UpdateFeatureFile updates the given section of the given feature file with the given new section.
func UpdateFeatureFile(filePath, oldSection, newSection string) error {
	// read file
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read feature file: %w", err)
	}

	// create lines
	fileLines := stringslice.Lines(string(content))
	// Normalize file lines for searching (replace SHA placeholders and actual SHAs with "SHA")
	normalizedFileLines := ReplaceSHAPlaceholder(fileLines)
	normalizedFileLines = ReplaceSHA(normalizedFileLines)
	normalizedFileLines = NormalizeTableWhitespace(normalizedFileLines)
	// Normalize old section for searching
	oldSectionLines := stringslice.TrimEmptyLines(stringslice.Lines(oldSection))
	oldSectionLines = ReplaceSHAPlaceholder(oldSectionLines)
	oldSectionLines = ReplaceSHA(oldSectionLines)
	oldSectionLines = NormalizeTableWhitespace(oldSectionLines)
	// Normalize new section (collapse excessive whitespace, preserve actual SHAs)
	newSectionLines := stringslice.TrimEmptyLines(stringslice.Lines(newSection))
	newSectionLines = NormalizeTableWhitespace(newSectionLines)

	// find the section in the file
	startLine, found := stringslice.LocateSection(normalizedFileLines, oldSectionLines)
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

// ReplaceSHAPlaceholder replaces all placeholders like "{{ sha.* }}" with "SHA".
func ReplaceSHAPlaceholder(lines []string) []string {
	return stringslice.ReplaceRegex(lines, regexp.MustCompile(`\{\{.*?\}\}`), "SHA")
}

func ReplaceSHA(lines []string) []string {
	return stringslice.ReplaceRegex(lines, regexp.MustCompile(`[a-z0-f]{40}`), "SHA")
}

// NormalizeTableWhitespace normalizes whitespace in Cucumber table rows by
// collapsing excessive whitespace (6+ consecutive spaces) while preserving normal column alignment.
func NormalizeTableWhitespace(lines []string) []string {
	excessiveSpaceRe := regexp.MustCompile(`\s{6,}`)
	result := make([]string, len(lines))
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "|") && strings.HasSuffix(trimmed, "|") {
			// This is a table row - collapse excessive whitespace to single space
			result[i] = excessiveSpaceRe.ReplaceAllString(trimmed, " ")
		} else {
			result[i] = trimmed
		}
	}
	return result
}
