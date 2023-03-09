package run

// DryRun implements the dry-run feature.
// The zero value is a non-activated DryRun.
type DryRun struct {
	active        bool   `exhaustruct:"optional"`
	currentBranch string `exhaustruct:"optional"`
}

// Activate enables dry-run.
func (dr *DryRun) Activate(currentBranch string) {
	dr.active = true
	dr.currentBranch = currentBranch
}

// ChangeBranch allows code to indicate to DryRun that the current branch has changed.
func (dr *DryRun) ChangeBranch(name string) {
	dr.currentBranch = name
}

// CurrentBranch provides the name of the current branch.
func (dr *DryRun) CurrentBranch() string {
	return dr.currentBranch
}

// IsActive indicates whether dry-run is active.
func (dr *DryRun) IsActive() bool {
	return dr.active
}
