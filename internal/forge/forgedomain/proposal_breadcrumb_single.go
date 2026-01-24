package forgedomain

import "strconv"

// ProposalBreadcrumbSingle indicates whether to add breadcrumbs to proposals of single branches,
// i.e. branches that are not part of a stack.
type ProposalBreadcrumbSingle bool

func (self ProposalBreadcrumbSingle) String() string {
	return strconv.FormatBool(bool(self))
}

func (self ProposalBreadcrumbSingle) Value() bool {
	return bool(self)
}
