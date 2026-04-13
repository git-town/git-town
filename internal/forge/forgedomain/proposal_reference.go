package forgedomain

// ProposalReferenceFallback returns a portable proposal reference when there is
// no forge-specific shorthand available.
func ProposalReferenceFallback(proposal ProposalData) string {
	title := proposal.Title.String()
	hasTitle := title != ""
	hasURL := proposal.URL != ""
	hasNumber := proposal.Number.Int() > 0
	switch {
	case hasTitle && hasURL:
		return "[" + title + "](" + proposal.URL + ")"
	case hasURL:
		return proposal.URL
	case hasTitle:
		return title
	case hasNumber:
		return "#" + proposal.Number.String()
	default:
		return ""
	}
}
