package flags

// type-safe access to the version CLI argument
func Version() (AddFunc, ReadBoolFlagFunc) {
	return Bool("version", "V", "Display the version number", FlagTypeNonPersistent)
}
