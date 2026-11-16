package gitdomain

type ConflictResolution string

func (self ConflictResolution) GitFlag() string {
	return "--" + self.String()
}

func (self ConflictResolution) String() string {
	return string(self)
}

const (
	ConflictResolutionOurs   ConflictResolution = "ours"
	ConflictResolutionTheirs ConflictResolution = "theirs"
)
