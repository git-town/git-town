package configdomain

// indicates whether a Git Town command should sync/display the lineage of a proposal
type Proposal bool

func (self Proposal) IsTrue() bool {
	return bool(self)
}
