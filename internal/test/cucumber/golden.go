package cucumber

import (
	"fmt"
	"os"
	"strings"
)

// updateFeatureFile updates the given section of the given feature file with the given new section.
func updateFeatureFile(filePath, oldSection, newSection string) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		panic(fmt.Sprintf("failed to read feature file: %v", err))
	}
	fileLines := strings.Split(string(content), "\n")
	oldSectionLines := trimmedLines(oldSection)
	newSectionLines := trimmedLines(newSection)

	// Find the section in the file
	startLine, found := locateSection(fileLines, oldSectionLines)
	if !found {
		fmt.Println("ERROR! Could not find section in feature file: ", filePath)
		fmt.Println("Expected section:\n", oldSection)
		return
	}

	// indent the new section
	indentation := getIndentation(fileLines[startLine])
	indentedNewTableLines := indentTableLines(newSectionLines, indentation)

	// Replace the old table lines with the new ones
	newLines := append([]string{}, fileLines[:startLine]...)
	newLines = append(newLines, indentedNewTableLines...)
	newLines = append(newLines, fileLines[startLine+len(oldSectionLines):]...)

	// Write back to the file
	newContent := strings.Join(newLines, "\n")
	//nolint:gosec // need permission 644 for feature files
	if err := os.WriteFile(filePath, []byte(newContent), 0o644); err != nil {
		panic(fmt.Sprintf("failed to write feature file: %v", err))
	}
}

// trimmedLines removes leading and trailing empty lines from a table string
func trimmedLines(tableStr string) []string {
	linesRaw := strings.Split(tableStr, "\n")
	// Filter out leading empty lines
	lines := make([]string, 0, len(linesRaw))
	for _, line := range linesRaw {
		if strings.TrimSpace(line) != "" || len(lines) > 0 {
			lines = append(lines, line)
		}
	}
	// Trim trailing empty lines
	for len(lines) > 0 && strings.TrimSpace(lines[len(lines)-1]) == "" {
		lines = lines[:len(lines)-1]
	}
	return lines
}

// locateSection locates a table in the file by matching unindented content
func locateSection(fileLines, tableLines []string) (int, bool) {
	for i := 0; i <= len(fileLines)-len(tableLines); i++ {
		if matchesTable(fileLines[i:], tableLines) {
			return i, true
		}
	}
	return -1, false
}

// matchesTable checks if the file lines match the table lines (ignoring indentation)
func matchesTable(fileLines, tableLines []string) bool {
	if len(fileLines) < len(tableLines) {
		return false
	}
	for j, tableLine := range tableLines {
		if strings.TrimSpace(fileLines[j]) != strings.TrimSpace(tableLine) {
			return false
		}
	}
	return true
}

// getIndentation extracts leading whitespace from a line
func getIndentation(line string) string {
	indentation := ""
	for _, ch := range line {
		if ch == ' ' || ch == '\t' {
			indentation += string(ch)
		} else {
			break
		}
	}
	return indentation
}

// indentTableLines applies indentation to each non-empty line
func indentTableLines(lines []string, indentation string) []string {
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
