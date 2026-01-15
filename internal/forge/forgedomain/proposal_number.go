package forgedomain

import "strconv"

// ProposalNumber is a number (ID) of a proposal.
//
// Example: https://github.com/git-town/git-town/pull/5977 has proposal number 5977.
type ProposalNumber int

func (self ProposalNumber) Int() int {
	return int(self)
}

func (self ProposalNumber) Int64() int64 {
	return int64(self)
}

func (self ProposalNumber) String() string {
	return strconv.Itoa(self.Int())
}

func NewProposalNumberFromFloat64(number float64) ProposalNumber {
	return ProposalNumber(int(number))
}
