package proposallineage

import (
	"strings"

	"github.com/git-town/git-town/v22/internal/git/gitdomain"
)

const (
	stackMarker = "<!-- branch-stack -->"
	stackStart  = "<!-- branch-stack-start -->"
	stackEnd    = "<!-- branch-stack-end -->"
)

func UpdateProposalBody(body gitdomain.ProposalBody, lineageSection string) gitdomain.ProposalBody {
	bodyStr := body.String()
	wrappedLineage := stackStart + "\n" + lineageSection + "\n" + stackEnd + "\n"

	// check for existing lineage section
	startIdx := strings.Index(bodyStr, stackStart)
	endIdx := strings.Index(bodyStr, stackEnd)
	if startIdx != -1 && endIdx != -1 {
		// Replace existing lineage section
		endLineEnd := endIdx + len(stackEnd)
		if endLineEnd < len(bodyStr) && bodyStr[endLineEnd] == '\n' {
			endLineEnd++
		}
		return gitdomain.ProposalBody(bodyStr[:startIdx] + wrappedLineage + bodyStr[endLineEnd:])
	}

	// here there is no lineage section, check for stack marker
	stackIdx := strings.Index(bodyStr, stackMarker)
	if stackIdx != -1 {
		// Insert lineage after the marker
		markerLineEnd := stackIdx + len(stackMarker)
		if markerLineEnd < len(bodyStr) && bodyStr[markerLineEnd] == '\n' {
			markerLineEnd++
		}
		return gitdomain.ProposalBody(bodyStr[:markerLineEnd] + wrappedLineage + bodyStr[markerLineEnd:])
	}

	// here there are no markers at all
	if bodyStr == "" {
		// empty body: just return the lineage section
		return gitdomain.ProposalBody(wrappedLineage)
	}

	// here the body is text without markers
	return gitdomain.ProposalBody(bodyStr + "\n\n" + wrappedLineage)
}
