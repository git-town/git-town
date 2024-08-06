package flags

// type-safe access to the "--all" command-line flag
func Stack(description string) (AddFunc, ReadBoolFlagFunc) {
	return Bool("stack", "s", description, FlagTypeNonPersistent)
}
