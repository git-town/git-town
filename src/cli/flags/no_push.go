package flags

// provides type-safe access to the "--no-push" command-line flag
func NoPush() (AddFunc, ReadBoolFlagFunc) {
	return Bool("no-push", "", "Do not push local branches", FlagTypePersistent)
}
