package flags

// type-safe access to the CLI arguments of type gitdomain.SwitchMerge
func SwitchMerge() (AddFunc, ReadBoolFlagFunc) {
	return Bool("merge", "m", "merge uncommitted changes into the target branch", FlagTypePersistent)
}
