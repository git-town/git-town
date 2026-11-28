package configdomain

// StageAll indicates whether a Git Town command should stage all open changes before committing.
type StageAll bool

func (self StageAll) Enabled() bool {
	return bool(self)
}
