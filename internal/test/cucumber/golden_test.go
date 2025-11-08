package cucumber_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/git-town/git-town/v22/internal/test/cucumber"
	"github.com/shoenig/test/must"
)

func TestUpdateFeatureFileWithCommands(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		file     string
		oldTable string
		newTable string
		want     string
	}{
		{
			name: "replace table with proper indentation",
			file: `
Feature: test
  Scenario: test
    Then Git Town runs the commands
      | BRANCH | COMMAND |
      | main   | git fetch |
    And some other step`[1:],
			oldTable: `
      | BRANCH | COMMAND |
      | main   | git fetch |`[1:],
			newTable: `
| BRANCH | COMMAND |
| main   | git pull |`[1:],
			want: `
Feature: test
  Scenario: test
    Then Git Town runs the commands
      | BRANCH | COMMAND |
      | main   | git pull |
    And some other step`[1:],
		},
		{
			name: "replace table with different number of rows",
			file: `
Feature: test
  Scenario: test
    Then Git Town runs the commands
      | BRANCH | COMMAND |
      | main   | git fetch |
    And done`[1:],
			oldTable: `
			| BRANCH | COMMAND |
			| main   | git fetch |`[1:],
			newTable: `
			| BRANCH | COMMAND |
			| main   | git init |
			| main   | git pull |`[1:],
			want: `
Feature: test
  Scenario: test
    Then Git Town runs the commands
      | BRANCH | COMMAND |
      | main   | git init |
      | main   | git pull |
    And done`[1:],
		},
		{
			name: "old table has more empty lines",
			file: `
Feature: test
  Scenario: test
    Then Git Town runs the commands
      | BRANCH | COMMAND |
      | main   | git fetch |
    And done`[1:],
			oldTable: `
			| BRANCH | COMMAND |
			| main   | git fetch |
`,
			newTable: `
			| BRANCH | COMMAND |
			| main   | git init |
			| main   | git pull |`[1:],
			want: `
Feature: test
  Scenario: test
    Then Git Town runs the commands
      | BRANCH | COMMAND |
      | main   | git init |
      | main   | git pull |
    And done`[1:],
		},
		{
			name: "newTable has more empty lines",
			file: `
Feature: test
  Scenario: test
    Then Git Town runs the commands
      | BRANCH | COMMAND |
      | main   | git fetch |
    And done`[1:],
			oldTable: `
			| BRANCH | COMMAND |
			| main   | git fetch |`[1:],
			newTable: `
			| BRANCH | COMMAND |
			| main   | git init |
			| main   | git pull |
`,
			want: `
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
			tmpDir := t.TempDir()
			tmpFile := filepath.Join(tmpDir, "test.feature")
			err := os.WriteFile(tmpFile, []byte(tt.file), 0o600)
			must.NoError(t, err)

			err = cucumber.UpdateFeatureFile(tmpFile, tt.oldTable, tt.newTable)
			must.NoError(t, err)

			result, err := os.ReadFile(tmpFile)
			must.NoError(t, err)
			must.EqOp(t, tt.want, string(result))
		})
	}
}
