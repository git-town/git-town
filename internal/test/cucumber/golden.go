package cucumber

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

// updateFeatureFileWithCommands updates the feature file with the actual commands table
func updateFeatureFileWithCommands(featureFilePath, oldTableStr, newTableStr string) {
	// Read the entire feature file
	content, err := os.ReadFile(featureFilePath)
	if err != nil {
		panic(fmt.Sprintf("failed to read feature file: %v", err))
	}

	// Split content into lines and prepare table lines
	fileLines := strings.Split(string(content), "\n")
	oldTableLines := trimTableLines(oldTableStr)
	newTableLines := trimTableLines(newTableStr)

	// Find the table in the file
	startLine, err := findTableInFile(fileLines, oldTableLines)
	if err != nil {
		fmt.Println("ERROR! Could not find expected table in feature file: ", featureFilePath)
		fmt.Println("Expected table:\n", oldTableStr)
		return
	}

	// Get indentation and apply it to the new table
	indentation := extractIndentation(fileLines[startLine])
	indentedNewTableLines := indentTableLines(newTableLines, indentation)

	// Replace the old table lines with the new ones
	newLines := append([]string{}, fileLines[:startLine]...)
	newLines = append(newLines, indentedNewTableLines...)
	newLines = append(newLines, fileLines[startLine+len(oldTableLines):]...)

	// Write back to the file
	newContent := strings.Join(newLines, "\n")
	//nolint:gosec // need permission 644 for feature files
	if err := os.WriteFile(featureFilePath, []byte(newContent), 0o644); err != nil {
		panic(fmt.Sprintf("failed to write feature file: %v", err))
	}
}

// trimTableLines removes leading and trailing empty lines from a table string
func trimTableLines(tableStr string) []string {
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

// findTableInFile locates a table in the file by matching unindented content
func findTableInFile(fileLines, tableLines []string) (int, error) {
	for i := 0; i <= len(fileLines)-len(tableLines); i++ {
		if matchesTable(fileLines[i:], tableLines) {
			return i, nil
		}
	}
	return -1, errors.New("could not find expected table in feature file")
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

// extractIndentation extracts leading whitespace from a line
func extractIndentation(line string) string {
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
