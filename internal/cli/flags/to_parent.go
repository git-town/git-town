package flags

// type-safe access to the "--to-parent" command-line flag
func ToParent() (AddFunc, ReadBoolFlagFunc) {
	return Bool("to-parent", "p", "allow shipping into non-perennial parent", FlagTypeNonPersistent)
}
