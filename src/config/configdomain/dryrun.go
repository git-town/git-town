package configdomain

// indicates whether a Git Town command should execute the commands or only display them
type DryRun bool

// TODO: remove
func (self DryRun) IsFalse() bool {
	return !bool(self)
}

func (self DryRun) IsTrue() bool {
	return bool(self)
}
