package flags

// type-safe access to the "--all" command-line flag
func All() (AddFunc, ReadBoolFlagFunc) {
	return Bool("all", "a", "sync all local branches", FlagTypeNonPersistent)
}
