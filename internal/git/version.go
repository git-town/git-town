package git

// the currently installed Git version
type Version struct {
	Major int
	Minor int
}

// indicates whether this version satisfies Git Town's minimum version requirement
func (self Version) IsAcceptableGitVersion() bool {
	return self.Major > 2 || (self.Major == 2 && self.Minor >= 30)
}

// provides an empty version instance
func emptyVersion() Version {
	var result Version
	return result
}
