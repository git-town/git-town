package forge

import (
	"strings"

	"github.com/git-town/git-town/v22/internal/git/gitdomain"
)

func ProposalBodyUpdateWithStackLineage(body gitdomain.ProposalBody, lineageContent string) gitdomain.ProposalBody {
	if lineageContent == "" {
		return body
	}
	const startMarker = "<!-- branch-stack -->"
	const endMarker = "<!-- branch-stack-end -->"

	bodyStr := body.String()

	// Create the full lineage section with both markers
	lineageSection := startMarker + "\n" + lineageContent + "\n" + endMarker

	// Find the start marker
	startIndex := strings.Index(bodyStr, startMarker)
	if startIndex != -1 {
		// Find where our section ends
		afterStart := bodyStr[startIndex:]

		var beforeSection, afterSection string
		beforeSection = bodyStr[:startIndex]
		// Look for the end marker
		endMarkerIndex := strings.Index(afterStart, endMarker)

		if endMarkerIndex != -1 {
			// End marker found - replace everything including the end marker
			afterSection = afterStart[endMarkerIndex+len(endMarker):]
		} else {
			// No end marker - preserve everything after our content
			// Find the end of the lineage content (look for double newline or end of string)
			contentAfterMarker := afterStart[len(startMarker):]

			// Try to find where the old lineage content ends
			// Look for the next section (typically starts with \n\n)
			doubleNewlineIndex := strings.Index(contentAfterMarker, "\n\n")
			if doubleNewlineIndex != -1 {
				afterSection = contentAfterMarker[doubleNewlineIndex:]
			} else {
				// No clear boundary found, set afterSection to everything after the startMarker
				afterSection = bodyStr[startIndex+len(startMarker):]
			}
		}

		return gitdomain.ProposalBody(beforeSection + lineageSection + afterSection)
	}

	// Marker doesn't exist - append it
	if bodyStr != "" {
		return gitdomain.ProposalBody(bodyStr + "\n\n" + lineageSection)
	}
	return gitdomain.ProposalBody(lineageSection)
}
