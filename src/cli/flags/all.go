package flags

// type-safe access to the "--all" command-line flag
func All() (AddFunc, ReadBoolFlagFunc) {
	return Bool("all", "a", "Sync all local branches", FlagTypeNonPersistent)
}
