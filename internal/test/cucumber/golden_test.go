package cucumber_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/git-town/git-town/v24/internal/test/cucumber"
	"github.com/shoenig/test/must"
)

func TestNormalizeWhitespace(t *testing.T) {
	t.Parallel()

	tests := []struct {
		give []string
		name string
		want []string
	}{
		{
			give: []string{
				"one    two",
				"three     four",
			},
			name: "multiple spaces",
			want: []string{
				"one two",
				"three four",
			},
		},
		{
			give: []string{
				"one\t\t\t\ttwo",
				"three\t\tfour",
			},
			name: "tabs",
			want: []string{
				"one two",
				"three four",
			},
		},
		{
			give: []string{
				"one \t  \t two",
				"three\t \tfour",
			},
			name: "mixed whitespace",
			want: []string{
				"one two",
				"three four",
			},
		},
		{
			give: []string{
				"one two",
				"three four five",
			},
			name: "single spaces unchanged",
			want: []string{
				"one two",
				"three four five",
			},
		},
		{
			give: []string{
				"  leading",
				"trailing  ",
				"  both  ",
			},
			name: "leading and trailing whitespace",
			want: []string{
				" leading",
				"trailing ",
				" both ",
			},
		},
		{
			give: []string{
				"nowhitespace",
				"stillnone",
			},
			name: "no whitespace",
			want: []string{
				"nowhitespace",
				"stillnone",
			},
		},
		{
			give: []string{
				"",
				"not empty",
				"",
			},
			name: "empty strings",
			want: []string{
				"",
				"not empty",
				"",
			},
		},
		{
			give: []string{},
			name: "empty slice",
			want: []string{},
		},
		{
			give: []string{
				"one\n\n\ntwo",
				"three\r\rfour",
				"five\f\fsix",
			},
			name: "newlines and other whitespace",
			want: []string{
				"one two",
				"three four",
				"five six",
			},
		},
		{
			give: []string{
				"      | BRANCH      | COMMAND   |",
				"      | main   | git fetch |",
			},
			name: "real-world example with table formatting",
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
		"gggggggggggggggggggggggggggggggggggggggg", // invalid hex (should not match)
		"0123456789abcdef0123456789abcdef01234567", // valid hex (should match)
	}
	want := []string{
		"SHA",
		"no sha",
		"",
		"gggggggggggggggggggggggggggggggggggggggg", // invalid hex (should not be replaced)
		"SHA",
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
		file     string
		name     string
		newTable string
		oldTable string
		want     string
	}{
		{
			file: `
Feature: test
  Scenario: test
    Then Git Town runs the commands
      | BRANCH | COMMAND |
      | main   | git fetch |
    And some other step`[1:],
			name: "replace table with proper indentation",
			newTable: `
| BRANCH | COMMAND |
| main   | git pull |`[1:],
			oldTable: `
      | BRANCH | COMMAND |
      | main   | git fetch |`[1:],
			want: `
Feature: test
  Scenario: test
    Then Git Town runs the commands
      | BRANCH | COMMAND |
      | main | git pull |
    And some other step`[1:],
		},
		{
			file: `
Feature: test
  Scenario: test
    Then Git Town runs the commands
      | BRANCH      | COMMAND   |
      | development | git fetch |
    And done`[1:],
			name: "replace table with different number of rows",
			newTable: `
			| BRANCH | COMMAND |
			| main   | git init |
			| main   | git pull |`[1:],
			oldTable: `
			| BRANCH      | COMMAND |
			| development | git fetch |`[1:],
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
			file: `
Feature: test
  Scenario: test
    Then Git Town runs the commands
      | BRANCH | COMMAND |
      | main   | git fetch |
    And done`[1:],
			name: "old table has more empty lines",
			newTable: `
			| BRANCH | COMMAND |
			| main | git init |
			| main | git pull |`[1:],
			oldTable: `

			| BRANCH | COMMAND |
			| main | git fetch |

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
			file: `
Feature: test
  Scenario: test
    Then Git Town runs the commands
      | BRANCH | COMMAND |
      | main   | git fetch |
    And done`[1:],
			name: "newTable has more empty lines",
			newTable: `

			| BRANCH | COMMAND |
			| main | git init |
			| main | git pull |

`,
			oldTable: `
			| BRANCH | COMMAND |
			| main | git fetch |`[1:],
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
			file: `
Feature: test
  Scenario: test
    Then Git Town runs the commands
      | BRANCH | COMMAND |
      | main   | git reset --hard {{ sha 'commit' }} |
    And done`[1:],
			name: "SHA placeholders",
			newTable: `
			| BRANCH | COMMAND                                                   |
			| main | git reset --soft d721118fcd545d37e87100b22ef13169160bdb3c |`[1:],
			oldTable: `
      | BRANCH | COMMAND                                                   |
      | main | git reset --hard d721118fcd545d37e87100b22ef13169160bdb3c |`[1:],
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
