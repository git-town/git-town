package forgedomain

import "strconv"

type ProposalNumber int

func (self ProposalNumber) Int() int {
	return int(self)
}

func (self ProposalNumber) String() string {
	return strconv.Itoa(int(self))
}

func NewProposalNumberFromFloat64(number float64) ProposalNumber {
	return ProposalNumber(int(number))
}
