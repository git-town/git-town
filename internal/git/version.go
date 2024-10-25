package git

// the currently installed Git version
type Version struct {
	Major int
	Minor int
}

// indicates whether the installed Git version supports the rebase.updateRefs config option
func (self Version) HasRebaseUpdateRefs() bool {
	return self.Major > 2 || (self.Major == 2 && self.Minor >= 38)
}

// indicates whether the installed Git version works for Git Town
func (self Version) IsMinimumRequiredGitVersion() bool {
	return self.Major > 2 || (self.Major == 2 && self.Minor >= 30)
}

// an empty version, to be used only in testing
func EmptyVersion() Version {
	return Version{} //exhaustruct:ignore
}
