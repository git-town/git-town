package cucumber_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/git-town/git-town/v22/internal/test/cucumber"
)

func TestIndentTableLines(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		lines       []string
		indentation string
		want        []string
	}{
		{
			name:        "add spaces",
			lines:       []string{"| A | B |", "| 1 | 2 |"},
			indentation: "    ",
			want:        []string{"    | A | B |", "    | 1 | 2 |"},
		},
		{
			name:        "add tabs",
			lines:       []string{"| A | B |", "| 1 | 2 |"},
			indentation: "\t\t",
			want:        []string{"\t\t| A | B |", "\t\t| 1 | 2 |"},
		},
		{
			name:        "no indentation",
			lines:       []string{"| A | B |", "| 1 | 2 |"},
			indentation: "",
			want:        []string{"| A | B |", "| 1 | 2 |"},
		},
		{
			name:        "preserve empty lines",
			lines:       []string{"| A | B |", "", "| 1 | 2 |"},
			indentation: "  ",
			want:        []string{"  | A | B |", "", "  | 1 | 2 |"},
		},
		{
			name:        "remove existing indentation and add new",
			lines:       []string{"  | A | B |", "    | 1 | 2 |"},
			indentation: "\t",
			want:        []string{"\t| A | B |", "\t| 1 | 2 |"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := cucumber.IndentSection(tt.lines, tt.indentation)
			if len(result) != len(tt.want) {
				t.Errorf("indentTableLines() returned %d lines, expected %d", len(result), len(tt.want))
				return
			}
			for i, line := range result {
				if line != tt.want[i] {
					t.Errorf("indentTableLines()[%d] = %q, expected %q", i, line, tt.want[i])
				}
			}
		})
	}
}

func TestUpdateFeatureFileWithCommands(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		initialContent string
		oldTable       string
		newTable       string
		wantResult     string
	}{
		{
			name: "replace table with proper indentation",
			initialContent: `
Feature: test
  Scenario: test
    Then Git Town runs the commands
      | BRANCH | COMMAND |
      | main   | git fetch |
    And some other step`[1:],
			oldTable: "| BRANCH | COMMAND |\n| main   | git fetch |",
			newTable: "| BRANCH | COMMAND |\n| main   | git pull |",
			wantResult: `
Feature: test
  Scenario: test
    Then Git Town runs the commands
      | BRANCH | COMMAND |
      | main   | git pull |
    And some other step`[1:],
		},
		{
			name: "replace table with different number of rows",
			initialContent: `
Feature: test
  Scenario: test
    Then Git Town runs the commands
      | BRANCH | COMMAND |
      | main   | git fetch |
    And done`[1:],
			oldTable: "| BRANCH | COMMAND |\n| main   | git fetch |",
			newTable: "| BRANCH | COMMAND |\n| main   | git init |\n| main   | git pull |",
			wantResult: `
Feature: test
  Scenario: test
    Then Git Town runs the commands
      | BRANCH | COMMAND |
      | main   | git init |
      | main   | git pull |
    And done`[1:],
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create temp file
			tmpDir := t.TempDir()
			tmpFile := filepath.Join(tmpDir, "test.feature")
			if err := os.WriteFile(tmpFile, []byte(tt.initialContent), 0o644); err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}

			// Run the function
			cucumber.UpdateFeatureFile(tmpFile, tt.oldTable, tt.newTable)

			// Read the result
			result, err := os.ReadFile(tmpFile)
			if err != nil {
				t.Fatalf("Failed to read result file: %v", err)
			}

			if string(result) != tt.wantResult {
				t.Errorf("updateFeatureFileWithCommands() result mismatch\nGot:\n%s\n\nExpected:\n%s", string(result), tt.wantResult)
			}
		})
	}
}
