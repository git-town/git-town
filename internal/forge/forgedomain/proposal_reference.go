package forgedomain

// ProposalReferenceFallback returns a portable proposal reference when there is
// no forge-specific shorthand available.
func ProposalReferenceFallback(proposal ProposalData) string {
	title := proposal.Title.String()
	switch {
	case title != "" && proposal.URL != "":
		return "[" + title + "](" + proposal.URL + ")"
	case proposal.URL != "":
		return proposal.URL
	case title != "":
		return title
	case proposal.Number.Int() > 0:
		return "#" + proposal.Number.String()
	default:
		return ""
	}
}
