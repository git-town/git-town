package forgedomain

import "strconv"

type ProposalBreadcrumbSingle bool

func (self ProposalBreadcrumbSingle) String() string {
	return strconv.FormatBool(bool(self))
}

func (self ProposalBreadcrumbSingle) Value() bool {
	return bool(self)
}
