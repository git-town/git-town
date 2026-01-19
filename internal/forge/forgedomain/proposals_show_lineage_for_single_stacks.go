package forgedomain

type ProposalsShowLineageSingleStack bool

func (self ProposalsShowLineageSingleStack) ShowLineage() bool {
	return bool(self)
}
