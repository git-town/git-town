package forgedomain

// ProposalInterface provides information about a change request on a forge.
// Alternative names are "pull request" or "merge request".
type ProposalInterface interface {
	Data() ProposalData
}

func CommitBody(data ProposalData, title string) string {
	result := title
	if body, has := data.Body.Get(); has {
		result += "\n\n"
		result += body.String()
	}
	return result
}
