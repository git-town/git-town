package flags

// SwitchMerge provides type-safe access to the CLI arguments of type gitdomain.ReadSwitchMergeFlagFunc.
func SwitchMerge() (AddFunc, ReadBoolFlagFunc) {
	return Bool("merge", "m", "merge uncommitted changes into the target branch", FlagTypePersistent)
}
