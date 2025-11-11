package cucumber_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/git-town/git-town/v22/internal/test/cucumber"
	"github.com/shoenig/test/must"
)

func TestNormalizeWhitespace(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		give []string
		want []string
	}{
		{
			name: "multiple spaces",
			give: []string{
				"one    two",
				"three     four",
			},
			want: []string{
				"one two",
				"three four",
			},
		},
		{
			name: "tabs",
			give: []string{
				"one\t\t\t\ttwo",
				"three\t\tfour",
			},
			want: []string{
				"one two",
				"three four",
			},
		},
		{
			name: "mixed whitespace",
			give: []string{
				"one \t  \t two",
				"three\t \tfour",
			},
			want: []string{
				"one two",
				"three four",
			},
		},
		{
			name: "single spaces unchanged",
			give: []string{
				"one two",
				"three four five",
			},
			want: []string{
				"one two",
				"three four five",
			},
		},
		{
			name: "leading and trailing whitespace",
			give: []string{
				"  leading",
				"trailing  ",
				"  both  ",
			},
			want: []string{
				" leading",
				"trailing ",
				" both ",
			},
		},
		{
			name: "no whitespace",
			give: []string{
				"nowhitespace",
				"stillnone",
			},
			want: []string{
				"nowhitespace",
				"stillnone",
			},
		},
		{
			name: "empty strings",
			give: []string{
				"",
				"not empty",
				"",
			},
			want: []string{
				"",
				"not empty",
				"",
			},
		},
		{
			name: "empty slice",
			give: []string{},
			want: []string{},
		},
		{
			name: "newlines and other whitespace",
			give: []string{
				"one\n\n\ntwo",
				"three\r\rfour",
				"five\f\fsix",
			},
			want: []string{
				"one two",
				"three four",
				"five six",
			},
		},
		{
			name: "real-world example with table formatting",
			give: []string{
				"      | BRANCH      | COMMAND   |",
				"      | main   | git fetch |",
			},
			want: []string{
				" | BRANCH | COMMAND |",
				" | main | git fetch |",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			have := cucumber.NormalizeWhitespace(tt.give)
			must.Eq(t, tt.want, have)
		})
	}
}

func TestReplaceSHA(t *testing.T) {
	t.Parallel()
	give := []string{
		"d721118fcd545d37e87100b22ef13169160bdb3c",
		"no sha",
		"",
	}
	want := []string{
		"SHA",
		"no sha",
		"",
	}
	have := cucumber.ReplaceSHA(give)
	must.Eq(t, want, have)
}

func TestReplaceSHAPlaceholder(t *testing.T) {
	t.Parallel()

	give := []string{
		"one {{ sha 'foo' }} two",
		"one {{ sha-in-origin 'bar' }} two",
		"git reset --hard {{ sha-initial 'alpha commit' }}",
		"no placeholder",
		"",
	}
	want := []string{
		"one SHA two",
		"one SHA two",
		"git reset --hard SHA",
		"no placeholder",
		"",
	}
	have := cucumber.ReplaceSHAPlaceholder(give)
	must.Eq(t, want, have)
}

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
      | main | git pull |
    And some other step`[1:],
		},
		{
			name: "replace table with different number of rows",
			file: `
Feature: test
  Scenario: test
    Then Git Town runs the commands
      | BRANCH      | COMMAND   |
      | development | git fetch |
    And done`[1:],
			oldTable: `
			| BRANCH      | COMMAND |
			| development | git fetch |`[1:],
			newTable: `
			| BRANCH | COMMAND |
			| main   | git init |
			| main   | git pull |`[1:],
			want: `
Feature: test
  Scenario: test
    Then Git Town runs the commands
      | BRANCH | COMMAND |
      | main | git init |
      | main | git pull |
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
			| main | git fetch |

`,
			newTable: `
			| BRANCH | COMMAND |
			| main | git init |
			| main | git pull |`[1:],
			want: `
Feature: test
  Scenario: test
    Then Git Town runs the commands
      | BRANCH | COMMAND |
      | main | git init |
      | main | git pull |
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
			| main | git fetch |`[1:],
			newTable: `

			| BRANCH | COMMAND |
			| main | git init |
			| main | git pull |

`,
			want: `
Feature: test
  Scenario: test
    Then Git Town runs the commands
      | BRANCH | COMMAND |
      | main | git init |
      | main | git pull |
    And done`[1:],
		},
		{
			name: "SHA placeholders",
			file: `
Feature: test
  Scenario: test
    Then Git Town runs the commands
      | BRANCH | COMMAND |
      | main   | git reset --hard {{ sha 'commit' }} |
    And done`[1:],
			oldTable: `
      | BRANCH | COMMAND                                                   |
      | main | git reset --hard d721118fcd545d37e87100b22ef13169160bdb3c |`[1:],
			newTable: `
			| BRANCH | COMMAND                                                   |
			| main | git reset --soft d721118fcd545d37e87100b22ef13169160bdb3c |`[1:],
			want: `
Feature: test
  Scenario: test
    Then Git Town runs the commands
      | BRANCH | COMMAND |
      | main | git reset --soft d721118fcd545d37e87100b22ef13169160bdb3c |
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

			err = cucumber.ChangeFeatureFile(tmpFile, tt.oldTable, tt.newTable)
			must.NoError(t, err)

			result, err := os.ReadFile(tmpFile)
			must.NoError(t, err)
			must.EqOp(t, tt.want, string(result))
		})
	}
}
