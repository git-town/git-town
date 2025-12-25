package gitdomain

type ConflictResolution string

const (
	ConflictResolutionOurs   ConflictResolution = "ours"
	ConflictResolutionTheirs ConflictResolution = "theirs"
)

func (self ConflictResolution) GitFlag() string {
	return "--" + self.String()
}

func (self ConflictResolution) String() string {
	return string(self)
}
