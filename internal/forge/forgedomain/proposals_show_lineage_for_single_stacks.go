package forgedomain

import "strconv"

type ProposalsShowLineageSingleStack bool

func (self ProposalsShowLineageSingleStack) Value() bool {
	return bool(self)
}

func (self ProposalsShowLineageSingleStack) String() string {
	return strconv.FormatBool(bool(self))
}
