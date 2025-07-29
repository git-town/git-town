package opcodes

import (
	"testing"

	"github.com/shoenig/test/must"
)

func TestUpdateBodyWithLineage(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		currentBody    string
		lineageContent string
		expected       string
	}{
		{
			name:           "empty body",
			currentBody:    "",
			lineageContent: "### This proposal is part of stack\n\n #123 [Title](URL)",
			expected:       "<!-- branch-stack -->\n### This proposal is part of stack\n\n #123 [Title](URL)\n<!-- branch-stack-end -->",
		},
		{
			name:           "body without marker",
			currentBody:    "This is the PR description",
			lineageContent: "### This proposal is part of stack\n\n #123 [Title](URL)",
			expected:       "This is the PR description\n\n<!-- branch-stack -->\n### This proposal is part of stack\n\n #123 [Title](URL)\n<!-- branch-stack-end -->",
		},
		{
			name:           "body with existing marker and end marker",
			currentBody:    "Description\n\n<!-- branch-stack -->\nOld lineage\n<!-- branch-stack-end -->\n\nMore content",
			lineageContent: "### New lineage",
			expected:       "Description\n\n<!-- branch-stack -->\n### New lineage\n<!-- branch-stack-end -->\n\nMore content",
		},
		{
			name:           "body with marker but missing end marker",
			currentBody:    "Description\n\n<!-- branch-stack -->\nOld lineage\n\nMore content",
			lineageContent: "### New lineage",
			expected:       "Description\n\n<!-- branch-stack -->\n### New lineage\n<!-- branch-stack-end -->\n\nMore content",
		},
		{
			name:           "multiple runs produce same result (idempotent)",
			currentBody:    "Initial content",
			lineageContent: "### Stack info",
			expected:       "Initial content\n\n<!-- branch-stack -->\n### Stack info\n<!-- branch-stack-end -->",
		},
		{
			name:           "body with marker at the end",
			currentBody:    "Some content\n\n<!-- branch-stack -->\nOld stack",
			lineageContent: "### New stack",
			expected:       "Some content\n\n<!-- branch-stack -->\n### New stack\n<!-- branch-stack-end -->",
		},
		{
			name:           "body with multiple HTML comments",
			currentBody:    "<!-- other comment -->\nContent\n<!-- branch-stack -->\nOld\n<!-- branch-stack-end -->\n<!-- another comment -->",
			lineageContent: "### Updated",
			expected:       "<!-- other comment -->\nContent\n<!-- branch-stack -->\n### Updated\n<!-- branch-stack-end -->\n<!-- another comment -->",
		},
		{
			name:           "lineage with special characters",
			currentBody:    "Description",
			lineageContent: "### Stack\n\n↳ #feature-1\n #123 [Title with [brackets]](URL) :point_left:",
			expected:       "Description\n\n<!-- branch-stack -->\n### Stack\n\n↳ #feature-1\n #123 [Title with [brackets]](URL) :point_left:\n<!-- branch-stack-end -->",
		},
		{
			name:           "preserve whitespace and formatting",
			currentBody:    "# Main Title\n\n## Description\n\nSome text",
			lineageContent: "### This proposal is part of stack\n\n    #123 [Indented](URL)\n  #124 [Also indented](URL2)",
			expected:       "# Main Title\n\n## Description\n\nSome text\n\n<!-- branch-stack -->\n### This proposal is part of stack\n\n    #123 [Indented](URL)\n  #124 [Also indented](URL2)\n<!-- branch-stack-end -->",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := updateProposalBodyWithStackLineage(tt.currentBody, tt.lineageContent)
			must.EqOp(t, tt.expected, result)

			// Test idempotency - running again should produce the same result
			if tt.name == "multiple runs produce same result (idempotent)" {
				secondRun := updateProposalBodyWithStackLineage(result, tt.lineageContent)
				must.EqOp(t, result, secondRun)
			}
		})
	}
}
