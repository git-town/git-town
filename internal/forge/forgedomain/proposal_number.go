package forgedomain

import (
	"errors"

	"github.com/git-town/git-town/v21/internal/messages"
)

type ProposalNumber int

func (self ProposalNumber) Validate() error {
	if self < 0 {
		return errors.New(messages.ProposalNoNumberGiven)
	}

	return nil
}

func (self ProposalNumber) ToInt() int {
	return int(self)
}
