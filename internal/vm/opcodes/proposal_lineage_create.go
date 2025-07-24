package opcodes

import (
	"fmt"
	"strings"

	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/vm/shared"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

type ProposalLineageCreate struct {
	Branch                  gitdomain.LocalBranchName
	ProposalLineageIn       configdomain.ProposalLineageIn
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *ProposalLineageCreate) Run(args shared.RunArgs) error {
	connector, hasConnector := args.Connector.Get()
	if !hasConnector {
		return forgedomain.UnsupportedServiceError()
	}
	findProposalFn, hasFindProposalFn := connector.FindProposalFn().Get()
	if !hasFindProposalFn {
		return fmt.Errorf("connector does not support finding proposals")
	}

	targetBranch := args.Config.Value.NormalConfig.Lineage.Parent(self.Branch)
	if targetBranch.IsNone() {
		return fmt.Errorf("current branch has no parent. Cannot find proposal to create lineage")
	}

	proposalData, err := findProposalFn(self.Branch, targetBranch.GetOrPanic())
	if err != nil {
		return err
	}

	if proposalData.IsNone() {
		return fmt.Errorf("current branch has no proposal")
	}

	builder := configdomain.NewProposalLineageBuilder(connector, args.Config.Value.MainAndPerennials()...)
	lineageInformation := args.Config.Value.NormalConfig.Lineage.BranchLineage(self.Branch)
	for _, curr := range lineageInformation {
		currParent := args.Config.Value.NormalConfig.Lineage.Parent(curr)
		builder.AddBranch(curr, currParent)
	}

	lineageAsString := builder.Build(self.Branch, self.ProposalLineageIn)

	switch self.ProposalLineageIn {
	case configdomain.ProposalLineageInTerminal:
		fmt.Print(lineageAsString)
		return nil
	case configdomain.ProposalLineageOperationInProposalBody:
		proposalDataUnwrapped := proposalData.GetOrPanic().Data
		currentBody := proposalDataUnwrapped.Data().Body.GetOrDefault()
		lineageString := lineageAsString.GetOrDefault()

		// Update body with lineage using our marker-based approach
		updatedBody := updateBodyWithLineage(currentBody, lineageString)

		op := &ProposalUpdateBody{
			Proposal:    proposalDataUnwrapped,
			UpdatedBody: Some(updatedBody),
		}
		return op.Run(args)
	case configdomain.ProposalLineageOperationInProposalComment:
		// TODO: Implement soon
	default:
		return nil
	}
	return nil
}

func updateBodyWithLineage(currentBody, lineageContent string) string {
	startMarker := "<!-- branch-stack -->"
	endMarker := "<!-- branch-stack-end -->"

	// Create the full lineage section with both markers
	lineageSection := startMarker + "\n" + lineageContent + "\n" + endMarker

	// Find the start marker
	startIndex := strings.Index(currentBody, startMarker)
	if startIndex != -1 {
		// Find where our section ends
		afterStart := currentBody[startIndex:]

		// Look for the end marker
		endMarkerIndex := strings.Index(afterStart, endMarker)

		var beforeSection, afterSection string
		beforeSection = currentBody[:startIndex]

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
				// No clear boundary found, replace everything
				afterSection = ""
			}
		}

		return beforeSection + lineageSection + afterSection
	}

	// Marker doesn't exist - append it
	if currentBody != "" {
		return currentBody + "\n\n" + lineageSection
	}
	return lineageSection
}
