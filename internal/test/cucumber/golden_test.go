package cucumber

import (
	"os"
	"path/filepath"
	"testing"
)

func TestTrimTableLines(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		give string
		want []string
	}{
		{
			name: "no empty lines",
			give: "| A | B |\n| 1 | 2 |",
			want: []string{"| A | B |", "| 1 | 2 |"},
		},
		{
			name: "trailing empty line",
			give: "| A | B |\n| 1 | 2 |\n",
			want: []string{"| A | B |", "| 1 | 2 |"},
		},
		{
			name: "multiple trailing empty lines",
			give: "| A | B |\n| 1 | 2 |\n\n\n",
			want: []string{"| A | B |", "| 1 | 2 |"},
		},
		{
			name: "leading empty lines",
			give: "\n\n| A | B |\n| 1 | 2 |",
			want: []string{"| A | B |", "| 1 | 2 |"},
		},
		{
			name: "empty string",
			give: "",
			want: []string{},
		},
		{
			name: "only empty lines",
			give: "\n\n\n",
			want: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := trimTableLines(tt.give)
			if len(result) != len(tt.want) {
				t.Errorf("trimTableLines() returned %d lines, expected %d", len(result), len(tt.want))
				return
			}
			for i, line := range result {
				if line != tt.want[i] {
					t.Errorf("trimTableLines()[%d] = %q, expected %q", i, line, tt.want[i])
				}
			}
		})
	}
}

func TestExtractIndentation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		give string
		want string
	}{
		{
			name: "spaces",
			give: "    | A | B |",
			want: "    ",
		},
		{
			name: "tabs",
			give: "\t\t| A | B |",
			want: "\t\t",
		},
		{
			name: "mixed spaces and tabs",
			give: "  \t  | A | B |",
			want: "  \t  ",
		},
		{
			name: "no indentation",
			give: "| A | B |",
			want: "",
		},
		{
			name: "empty line",
			give: "",
			want: "",
		},
		{
			name: "only whitespace",
			give: "    ",
			want: "    ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := extractIndentation(tt.give)
			if result != tt.want {
				t.Errorf("extractIndentation() = %q, expected %q", result, tt.want)
			}
		})
	}
}

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
			result := indentTableLines(tt.lines, tt.indentation)
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

func TestMatchesTable(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		fileLines  []string
		tableLines []string
		want       bool
	}{
		{
			name:       "exact match",
			fileLines:  []string{"| A | B |", "| 1 | 2 |"},
			tableLines: []string{"| A | B |", "| 1 | 2 |"},
			want:       true,
		},
		{
			name:       "match with different indentation",
			fileLines:  []string{"    | A | B |", "    | 1 | 2 |"},
			tableLines: []string{"| A | B |", "| 1 | 2 |"},
			want:       true,
		},
		{
			name:       "no match - different content",
			fileLines:  []string{"| A | B |", "| 3 | 4 |"},
			tableLines: []string{"| A | B |", "| 1 | 2 |"},
			want:       false,
		},
		{
			name:       "no match - file too short",
			fileLines:  []string{"| A | B |"},
			tableLines: []string{"| A | B |", "| 1 | 2 |"},
			want:       false,
		},
		{
			name:       "match - file has extra lines",
			fileLines:  []string{"| A | B |", "| 1 | 2 |", "| 3 | 4 |"},
			tableLines: []string{"| A | B |", "| 1 | 2 |"},
			want:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := matchesTable(tt.fileLines, tt.tableLines)
			if result != tt.want {
				t.Errorf("matchesTable() = %v, expected %v", result, tt.want)
			}
		})
	}
}

func TestFindTableInFile(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		fileLines  []string
		tableLines []string
		wantIdx    int
		wantError  bool
	}{
		{
			name: "table at beginning",
			fileLines: []string{
				"| A | B |",
				"| 1 | 2 |",
				"Some text",
			},
			tableLines: []string{"| A | B |", "| 1 | 2 |"},
			wantIdx:    0,
			wantError:  false,
		},
		{
			name: "table in middle",
			fileLines: []string{
				"Some text",
				"    | A | B |",
				"    | 1 | 2 |",
				"More text",
			},
			tableLines: []string{"| A | B |", "| 1 | 2 |"},
			wantIdx:    1,
			wantError:  false,
		},
		{
			name: "table at end",
			fileLines: []string{
				"Some text",
				"More text",
				"  | A | B |",
				"  | 1 | 2 |",
			},
			tableLines: []string{"| A | B |", "| 1 | 2 |"},
			wantIdx:    2,
			wantError:  false,
		},
		{
			name: "table not found",
			fileLines: []string{
				"| A | B |",
				"| 3 | 4 |",
			},
			tableLines: []string{"| A | B |", "| 1 | 2 |"},
			wantIdx:    -1,
			wantError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			idx, err := findTableInFile(tt.fileLines, tt.tableLines)
			if tt.wantError {
				if err == nil {
					t.Errorf("findTableInFile() expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("findTableInFile() unexpected error: %v", err)
				}
				if idx != tt.wantIdx {
					t.Errorf("findTableInFile() = %d, expected %d", idx, tt.wantIdx)
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
			updateFeatureFileWithCommands(tmpFile, tt.oldTable, tt.newTable)

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
