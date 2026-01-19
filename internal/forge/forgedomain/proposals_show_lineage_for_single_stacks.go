package forgedomain

import "strconv"

type ProposalsShowLineageSingleStack bool

func (self ProposalsShowLineageSingleStack) ShowLineage() bool {
	return bool(self)
}

func (self ProposalsShowLineageSingleStack) String() string {
	return strconv.FormatBool(bool(self))
}
