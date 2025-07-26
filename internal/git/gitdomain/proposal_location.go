package gitdomain

import . "github.com/git-town/git-town/v21/pkg/prelude"

type ProposalLocation struct {
	Id  int
	Url string
}

func NewProposalLocation(id int, url string) ProposalLocation {
	if !isValidProposalLocation(id, url) {
		panic("ProposalLocation must have a valid id and url")
	}
	return ProposalLocation{
		Id:  id,
		Url: url,
	}
}

func NewProposalLocationOption(id int, url string) Option[ProposalLocation] {
	if isValidProposalLocation(id, url) {
		return Some(NewProposalLocation(id, url))
	}
	return None[ProposalLocation]()
}

func isValidProposalLocation(id int, url string) bool {
	return id > 0 && len(url) > 0
}
